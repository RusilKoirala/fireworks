#ifndef Fireworks
  #define Fireworks "0.0.0"
#endif

[Setup]
AppId={{7B9E4C21-5A3F-4D8E-9C61-FIREWORKS0001}
AppName=Fireworks
AppVersion={#MyAppVersion}
AppPublisher=Rusil Koirala
DefaultDirName={localappdata}\Programs\Fireworks
DisableProgramGroupPage=yes
PrivilegesRequired=lowest
ChangesEnvironment=yes
OutputBaseFilename=fireworks-setup
Compression=lzma2
SolidCompression=yes
ArchitecturesInstallIn64BitMode=x64compatible

[Files]
Source: "fireworks-windows-amd64.exe"; DestDir: "{app}"; DestName: "fireworks.exe"; Flags: ignoreversion

[Registry]
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "Path"; \
  ValueData: "{olddata};{app}"; Check: NeedsAddPath(ExpandConstant('{app}'))

[Code]
function NeedsAddPath(Dir: string): Boolean;
var
  Paths: string;
begin
  if not RegQueryStringValue(HKCU, 'Environment', 'Path', Paths) then
  begin
    Result := True;
    exit;
  end;
  Result := Pos(';' + Uppercase(Dir) + ';', ';' + Uppercase(Paths) + ';') = 0;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
var
  Paths, Dir: string;
  P: Integer;
begin
  if CurUninstallStep = usPostUninstall then
  begin
    Dir := ExpandConstant('{app}');
    if RegQueryStringValue(HKCU, 'Environment', 'Path', Paths) then
    begin
      P := Pos(';' + Uppercase(Dir), Uppercase(Paths));
      if P > 0 then
      begin
        Delete(Paths, P, Length(Dir) + 1);
        RegWriteExpandStringValue(HKCU, 'Environment', 'Path', Paths);
      end;
    end;
  end;
end;