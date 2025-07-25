### Val Launcher

Not ready for use

## This is NOT a hack client

Val Launcher adds in some simple configuration to automaticly set the main menu background image

## How To Use

Specify the filepath of the game as well as the input mp4s and output location in the config.json
Multiple inputs will cause the launcher to choose one at random

```json config.json
{
  "exe_filepath": "C:/Riot Games/Riot Client/RiotClientServices.exe",
  "changes": [
    {
      "description": "homescreen",
      "inputs": [
        "C:/Users/henry/Projects/Val_Launcher/resources/red_dress_1.mp4",
        "C:/Users/henry/Projects/Val_Launcher/resources/red_dress_2.mp4",
        "C:/Users/henry/Projects/Val_Launcher/resources/pink_hair.mp4",
        "C:/Users/henry/Projects/Val_Launcher/resources/black_dress.mp4"
      ],
      "ouput": "C:/Riot Games/VALORANT/live/ShooterGame/Content/Movies/Menu/11_00_Homescreen.mp4"
    }
  ]
}
```

## Bannable?

Some people have reported getting banned for changing the background file, use at your own risk

## Plans

A gui is planned for selecting images and the files they replace
