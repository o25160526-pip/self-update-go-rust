package updater

import (
    "encoding/json"
    "os"
    "path/filepath"
    "strings"
)

const defaultManifestURL = "https://github.com/o25160526-pip/self-update-go-rust/releases/latest/download/manifest-go.json"
var ManifestURL = defaultManifestURL

type UpdateServerConfig struct { ActiveHost string `json:"activeHost"`; Hosts map[string]UpdateServerHost `json:"hosts"` }
type UpdateServerHost struct { BaseURL string `json:"baseUrl"`; GoManifest string `json:"goManifest"`; RustManifest string `json:"rustManifest"` }
func (h UpdateServerHost) GoURL() string { return strings.TrimRight(h.BaseURL, "/") + "/" + h.GoManifest }
func LoadUpdateServerConfig(path string) (UpdateServerConfig, error) { b, err := os.ReadFile(path); if err != nil { return UpdateServerConfig{}, err }; var c UpdateServerConfig; if err := json.Unmarshal(b, &c); err != nil { return UpdateServerConfig{}, err }; return c, nil }
func applyUpdateServerConfig() { var candidates []string; if p:=os.Getenv("GO_DEMO_UPDATE_SERVER_CONFIG"); p!="" { candidates=append(candidates,p) }; if e,err:=os.Executable(); err==nil { candidates=append(candidates,filepath.Join(filepath.Dir(e),"update-server.json")) }; if w,err:=os.Getwd(); err==nil { candidates=append(candidates,filepath.Join(w,"update-server.json")) }; for _,p:=range candidates { c,err:=LoadUpdateServerConfig(p); if err!=nil { continue }; h,ok:=c.Hosts[c.ActiveHost]; if ok && h.BaseURL!="" && h.GoManifest!="" { ManifestURL=h.GoURL(); return } } }
func init(){ applyUpdateServerConfig() }
func DataDir()(string,error){ base,err:=os.UserCacheDir(); if err!=nil{return os.MkdirTemp("","go-demo")}; return filepath.Join(base,"go-demo"),nil }
func StatePath(dir string)string{return filepath.Join(dir,"state.json")}
func DownloadDir(dir string)string{return filepath.Join(dir,"update-tmp")}
func HealthDir(dir string)string{return filepath.Join(dir,"health")}
