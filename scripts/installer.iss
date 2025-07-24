; MyAppInstaller.iss
[Setup]
AppName=Val_Launcher
AppVersion=1.0
DefaultDirName={pf}\Val_Launcher
DefaultGroupName=Val_Launcher
OutputDir=release
OutputBaseFilename=Val_Launcher_Wizard
Compression=lzma
SolidCompression=yes

[Files]
Source: "release\launch.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "release\config.json"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\MyApp"; Filename: "{app}\MyApp.exe"
Name: "{group}\Uninstall MyApp"; Filename: "{uninstallexe}"
