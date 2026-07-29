package updater

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"strings"
	"time"
)

type UpdateResult struct { HasUpdate bool `json:"hasUpdate"`; LatestVersion string `json:"latestVersion,omitempty"`; ReleaseNotes string `json:"releaseNotes,omitempty"`; URL string `json:"url,omitempty"`; Manifest *Manifest `json:"manifest,omitempty"`; Source string `json:"source,omitempty"` }
type Checker struct { Client *http.Client }
func NewChecker() *Checker { return &Checker{Client:&http.Client{Timeout:10*time.Minute}} }

func (c *Checker) Check(ctx context.Context, manifestURL, currentVersion string, policy UpdatePolicy) (*UpdateResult,error) { return c.checkOne(ctx, UpdateSource{Name:"configured",ManifestURL:manifestURL,Enabled:true}, currentVersion, policy) }

// CheckSources fetches every enabled source, ignores unavailable sources, and chooses the highest valid version.
// GitHub is normally priority 0, so it wins ties and is preferred for exact rollback lookups.
func (c *Checker) CheckSources(ctx context.Context, sources []UpdateSource, currentVersion string, policy UpdatePolicy) (*UpdateResult,error) {
	if len(sources)==0 { sources=UpdateSources }
	results:=make([]*UpdateResult,0,len(sources)); var errs []string
	for _, source := range sources { r,err:=c.checkOne(ctx,source,currentVersion,policy); if err!=nil { errs=append(errs,source.Name+": "+err.Error()); continue }; if r.HasUpdate { results=append(results,r) } }
	if len(results)==0 { if len(errs)==len(sources) { return nil,fmt.Errorf("all update sources failed: %s",strings.Join(errs,"; ")) }; return &UpdateResult{},nil }
	sort.SliceStable(results,func(i,j int)bool{return results[i].LatestVersion>results[j].LatestVersion})
	return results[len(results)-1],nil
}
func (c *Checker) checkOne(ctx context.Context, source UpdateSource, currentVersion string, policy UpdatePolicy)(*UpdateResult,error){ m,err:=FetchManifest(ctx,c.Client,source.ManifestURL);if err!=nil{return nil,err};if err=m.Validate();err!=nil{return nil,err};if m.Channel!=policy.Channel||m.Platform!=PlatformString(){return &UpdateResult{},nil};newer,err:=IsNewer(currentVersion,m.Version);if err!=nil{return nil,fmt.Errorf("compare version: %w",err)};if !newer&&!policy.AllowDowngrade{return &UpdateResult{},nil};return &UpdateResult{HasUpdate:true,LatestVersion:m.Version,ReleaseNotes:ParseReleaseNotes(m.ReleaseNotes),URL:m.URL,Manifest:m,Source:source.Name},nil}

// FindExactVersion checks GitHub first, then all configured sources, for rollback or recovery.
func (c *Checker) FindExactVersion(ctx context.Context, version string, sources []UpdateSource, policy UpdatePolicy)(*Manifest,string,error){ if len(sources)==0{sources=UpdateSources};sort.SliceStable(sources,func(i,j int)bool{if sources[i].Kind=="github"&&sources[j].Kind!="github"{return true};return sources[i].Priority<sources[j].Priority});var errs []string;for _,s:=range sources{m,err:=FetchManifest(ctx,c.Client,s.ManifestURL);if err!=nil{errs=append(errs,s.Name+": "+err.Error());continue};if m.Version==version&&m.Validate()==nil{return m,s.Name,nil}};return nil,"",fmt.Errorf("version %s not found in configured sources: %s",version,strings.Join(errs,"; "))}
func PlatformString()string{arch:=runtime.GOARCH;if arch=="amd64"{arch="x86_64"};return fmt.Sprintf("%s-%s",runtime.GOOS,arch)}
type UpdateState string
const(StateChecking UpdateState="checking";StateUpToDate UpdateState="up-to-date";StateUpdateAvailable UpdateState="update-available";StateDownloading UpdateState="downloading";StateVerifying UpdateState="verifying";StateInstalling UpdateState="installing";StateRestarting UpdateState="restarting";StateFailed UpdateState="failed";StateRolledBack UpdateState="rolled-back")
func IsValidState(s string)bool{switch UpdateState(s){case StateChecking,StateUpToDate,StateUpdateAvailable,StateDownloading,StateVerifying,StateInstalling,StateRestarting,StateFailed,StateRolledBack:return true};return false}
func MockCheckForUpdate(currentVersion string,hasUpdate bool)*UpdateResult{if !hasUpdate{return &UpdateResult{}};return &UpdateResult{HasUpdate:true,LatestVersion:"1.0.1",ReleaseNotes:"Mock: Fix update flow and restart",URL:"mock://artifact-url"}}
func ParseReleaseNotes(notes string)string{notes=strings.TrimSpace(notes);if len(notes)>200{return notes[:197]+"..."};return notes}
