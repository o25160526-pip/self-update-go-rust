package main

import ("context"; "fmt"; "runtime"; "time"; "go-demo/internal/version"; "go-demo/updater")
type App struct { ctx context.Context; svc *updater.Service }
func NewApp(svc *updater.Service)*App{return &App{svc:svc}}
func (a *App) startup(ctx context.Context){a.ctx=ctx;if !a.svc.Policy.AutoCheckOnStartup{return};go func(){checkCtx,cancel:=context.WithTimeout(ctx,10*time.Minute);defer cancel();_,_=a.svc.CheckAndUpdate(checkCtx)}()}
func (a *App) GetVersion()string{return version.Version}
func (a *App) GetInfo()map[string]string{return map[string]string{"app":"go-demo","version":version.Version,"os":runtime.GOOS,"arch":runtime.GOARCH,"keyId":updater.PinnedKeyID,"endpoint":updater.ManifestURL}}
func (a *App) GetUpdateState()string{return a.svc.State()}
func (a *App) GetLogs()[]string{return a.svc.Logs()}
func (a *App) CheckForUpdate()string{ctx,cancel:=context.WithTimeout(a.ctx,5*time.Minute);defer cancel();_,_=a.svc.CheckAndUpdate(ctx);return a.svc.State()}
func (a *App) RollbackPrevious() string { if err:=a.svc.RollbackPrevious();err!=nil{return err.Error()};return "rollback-started" }
func (a *App) Greet(name string)string{return fmt.Sprintf("Hello %s, It's show time!",name)}
