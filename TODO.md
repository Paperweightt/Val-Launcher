# DO IT

## todo list

- [ ] setup installer
  - [ ] set iss cmd
  - [ ] set iss into script/build.ps1
- [ ] config ui
  - [ ] setup dependenciess
- [ ] setup dir input so collections of files can be input

## inputs

inputs: (dir|file)[]

pros: dont have to list out as many inputs in the config.json
cons: cant list out dir to replace dir, not sure if this works currently but might

```go launch/main.go
type Config struct {
	ExeFilepath string `json:"exe_filepath"`
	Changes     []struct {
		Description string   `json:"description"`
		Inputs      []string `json:"inputs"` // keep the same but handle differently
		Ouput       string   `json:"ouput"`
	} `json:"changes"`
}
```
