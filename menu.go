package main

import (
  "io/ioutil"
  "fmt"
  "os"
  "os/exec"
  "time"
  "encoding/xml"
  "sort"
  "bytes"
  "path/filepath"

  "github.com/veandco/go-sdl2/sdl"
  "github.com/veandco/go-sdl2/ttf"
  "gopkg.in/yaml.v2"

  //this makes our life easier by giving the sdl.Color corresponding to the colorname (converted from golang.org/x/image/colornames)
  "./sdlcolornames"
)

type (
  items struct {
    Desc string
    Desc2 string
    Cmd string
    Undocmd string
    Check string
  }
  menulist struct {
    Name string
    Options map[string]items
  }
)

var (
  winTitle string = "Game Boy Zero config menu"
  winWidth, winHeight int32
  fonttype string = "resources/opensans_hebrew_condensed_regular.ttf"
  fontsize int
  es_inputcfg string = "/opt/retropie/configs/all/emulationstation/es_input.cfg"
  backgroundcolor sdl.Color = sdlcolornames.Whitesmoke
  menucursorcolor sdl.Color = sdlcolornames.Grey
  positiontitle sdl.Rect = sdl.Rect{80, 5, 0, 0}
  positionmenu sdl.Rect //= sdl.Rect{30, 55, 0, 0}
  menucursor sdl.Rect
  buttonpressed string
)
func run() int {
  type (
    Es_input_button struct {
      Name  string `xml:"name,attr"`
      Type  string `xml:"type,attr"`
      Id    int `xml:"id,attr"`
      Value uint8 `xml:"value,attr"`
    }
    Es_input_config struct {
     Type string `xml:"type,attr"`
     Guid string `xml:"deviceGUID,attr"`
     Buttons []Es_input_button `xml:"input"`
    }
    Es_input struct {
      Config Es_input_config `xml:"inputConfig"`
    }
    Joystick_buttons struct {
      Type  string
      Id    int
      Value uint8
    }
  )

  var (
    window *sdl.Window
    font *ttf.Font
    surface *sdl.Surface
    input Es_input
    joystick_up, joystick_down, joystick_left, joystick_right, joystick_a, joystick_b Joystick_buttons
    menuitens map[int]string
    displaybounds sdl.Rect
  )

  xmlfile, _ := os.Open(es_inputcfg)
  es_inputfile, _ := ioutil.ReadAll(xmlfile)

  xml.Unmarshal([]byte(es_inputfile), &input)

  //we need to create a "map" of the buttons at emulationstation so we can check what the pressed button is
  for n:=0;n<len(input.Config.Buttons);n++ {
    switch input.Config.Buttons[n].Name {
    case "up":
      joystick_up.Type = input.Config.Buttons[n].Type
      joystick_up.Id = input.Config.Buttons[n].Id
      joystick_up.Value = input.Config.Buttons[n].Value
    case "down":
      joystick_down.Type = input.Config.Buttons[n].Type
      joystick_down.Id = input.Config.Buttons[n].Id
      joystick_down.Value = input.Config.Buttons[n].Value
    case "left":
      joystick_left.Type = input.Config.Buttons[n].Type
      joystick_left.Id = input.Config.Buttons[n].Id
      joystick_left.Value = input.Config.Buttons[n].Value
    case "right":
      joystick_right.Type = input.Config.Buttons[n].Type
      joystick_right.Id = input.Config.Buttons[n].Id
      joystick_right.Value = input.Config.Buttons[n].Value 
    case "a":
      joystick_a.Type = input.Config.Buttons[n].Type
      joystick_a.Id = input.Config.Buttons[n].Id
      joystick_a.Value = input.Config.Buttons[n].Value
    case "b":
      joystick_b.Type = input.Config.Buttons[n].Type
      joystick_b.Id = input.Config.Buttons[n].Id
      joystick_b.Value = input.Config.Buttons[n].Value
    }
  }

  sdl.Init(sdl.INIT_VIDEO | sdl.INIT_JOYSTICK)

  ttf.Init()

  sdl.GetDisplayBounds(0, &displaybounds)
  
  winWidth =  displaybounds.W
  winHeight = displaybounds.H

  window, _ = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN | sdl.WINDOW_FULLSCREEN)
  defer window.Destroy()

  menu := menulistyaml()

  //it's cooler/nicer if the menu is sorted, besides if we don't sort it, the menu order will change at every run
  var sortkey []string
  for key := range menu.Options {
    sortkey = append(sortkey, key)
  }
  sort.Strings(sortkey)

  //we need to have a map crossing position x itemkey
  n := 1
  menuitens = make(map[int]string)
  for _, k := range sortkey {
    menuitens[n] = k
    n++
  }
  
  fontsize = int(winWidth/20)

  //we need a font to work with =p for now we'll use same font/size for title and itens
  dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
  fontdir := dir + "/" + fonttype
  font, _ = ttf.OpenFont(fontdir, fontsize)
  
  //so we can set the title to be centralized :)
  titlewidth, _, _ := font.SizeUTF8(menu.Name)
  positiontitle.X = int32(int(winWidth)/2 - titlewidth/2)

  //we need to know how many itens fit at the screen size, it's all math: maxitems = (screenheight-titleheight)/lineheight
  titleheight := int(positiontitle.Y)+font.Height()
  menumaxitensatscreen := (int(winHeight)-titleheight)/font.Height()

  positionmenu = sdl.Rect{30, int32(titleheight+5), 0, 0}
  //so we can set the menucursor height to be relative to font size :)
  menucursor = sdl.Rect{10, positionmenu.Y,winWidth - 20,int32(font.Height())}
  
  //we need to know the active surface at screen so we can "paste" our texts there later
  surface, _ = window.GetSurface()

  //create the renderer that will colour our background and menu cursor
  renderer, _ := sdl.CreateSoftwareRenderer(surface)

  last := time.Now()
  running := true
  positionmenucursor := 1
  firstmenuitem := positionmenucursor
  titletext, menutext := drawmenuitens(menuitens, firstmenuitem, menumaxitensatscreen, menu, font, surface)
  for running {
    renderer.SetDrawColor(backgroundcolor.R, backgroundcolor.G, backgroundcolor.B, backgroundcolor.A)
    renderer.Clear()
    renderer.Present()

    drawmenucursor(renderer, menucursorcolor, menucursor)

    titletext.Blit(nil, surface, &positiontitle)
    menutext.Blit(nil, surface, &positionmenu)
    window.UpdateSurface()

    if time.Since(last).Seconds() > 0.15 {
      switch buttonpressed {
      case "up":
        if positionmenucursor > 1 {
          if positionmenucursor != firstmenuitem {
            menucursor.Y -= menucursor.H
            positionmenucursor--
          } else {
            firstmenuitem--
            positionmenucursor--
            titletext, menutext = drawmenuitens(menuitens, firstmenuitem, menumaxitensatscreen, menu, font, surface)
          }
          last = time.Now()
        }
      case "down":
        if positionmenucursor < len(menuitens) {
          //we do this so we can scroll the screen down when we reach the bottom and there's more itens to come
          //the "firstmenuitem + 1" is a desperate fix to the "screen scrolling down" when it wasn't supposed to
          if positionmenucursor - firstmenuitem + 1 < menumaxitensatscreen {
            menucursor.Y += menucursor.H
            positionmenucursor++
          } else {
            firstmenuitem++
            positionmenucursor++
            titletext, menutext = drawmenuitens(menuitens, firstmenuitem, menumaxitensatscreen, menu, font, surface)
          }
          last = time.Now()
        }
      case "a":
        oldtitle := menu.Name
        menu.Name = "Working..."
        //so we can set the title to be centralized :)
        titlewidth, _, _ = font.SizeUTF8(menu.Name)
        positiontitle.X = int32(int(winWidth)/2 - titlewidth/2)
        titletext, menutext = drawmenuitens(menuitens, firstmenuitem, menumaxitensatscreen, menu, font, surface)
        renderer.SetDrawColor(backgroundcolor.R, backgroundcolor.G, backgroundcolor.B, backgroundcolor.A)
        renderer.Clear()
        renderer.Present()
        drawmenucursor(renderer, sdlcolornames.Lightgrey, menucursor)
        titletext.Blit(nil, surface, &positiontitle)
        menutext.Blit(nil, surface, &positionmenu)
        window.UpdateSurface()
        if menu.Options[menuitens[positionmenucursor]].Check != "" {
          if err := exec.Command("/bin/bash", "-c", menu.Options[menuitens[positionmenucursor]].Check).Run(); err != nil {
            runcommand(menu.Options[menuitens[positionmenucursor]].Cmd)
          } else {
            runcommand(menu.Options[menuitens[positionmenucursor]].Undocmd)
          }
        } else {
          runcommand(menu.Options[menuitens[positionmenucursor]].Cmd)
        }
        menu.Name = oldtitle
        //so we can set the title to be centralized :)
        titlewidth, _, _ = font.SizeUTF8(menu.Name)
        positiontitle.X = int32(int(winWidth)/2 - titlewidth/2)
        last = time.Now()
        titletext, menutext = drawmenuitens(menuitens, firstmenuitem, menumaxitensatscreen, menu, font, surface)
      case "b":
        running = false
      }
    }
    event := sdl.PollEvent()
    switch t := event.(type) {
      case *sdl.QuitEvent:
        running = false
      case *sdl.KeyboardEvent:
        if input.Config.Type == "keyboard" {
          if t.Type == sdl.KEYDOWN {
            switch int(t.Keysym.Sym) {
            case joystick_up.Id:
              buttonpressed = "up"
            case joystick_down.Id:
              buttonpressed = "down"
            case joystick_a.Id:
              buttonpressed = "a"
            case joystick_b.Id:
              buttonpressed = "b"
            }
          }
          if t.Type == sdl.KEYUP {
            buttonpressed = ""
          }
        }

      case *sdl.JoyHatEvent:
        if input.Config.Type == "joystick" {
          switch t.Value {
            case joystick_up.Value:
              buttonpressed = "up"
            case joystick_down.Value:
              buttonpressed = "down"
            case 0:
              buttonpressed = ""
          }
        }
      case *sdl.JoyButtonEvent:
        if input.Config.Type == "joystick" {
          switch t.Button {
          case uint8(joystick_a.Id):
            if t.State == joystick_a.Value {
              buttonpressed = "a"
            } else {
              buttonpressed = ""
            }
          case uint8(joystick_b.Id):
            if t.State == joystick_a.Value {
              buttonpressed = "b"
            } else {
              buttonpressed = ""
            }
          }
        }
      case *sdl.JoyDeviceEvent:
        if t.Type == sdl.JOYDEVICEADDED {
          sdl.JoystickOpen(int(t.Which))
        }
      }
  }
  return 0
}

func menulistyaml() menulist{
  var menu menulist

  yamlfile, _ := os.Open(os.Args[1])
  menulistfile, _ := ioutil.ReadAll(yamlfile)

  yaml.Unmarshal([]byte(menulistfile),&menu)

  return menu
}

func drawmenucursor(renderer *sdl.Renderer, menucursorcolor sdl.Color, rect sdl.Rect) {
  renderer.SetDrawColor(menucursorcolor.R, menucursorcolor.G, menucursorcolor.B, menucursorcolor.A)
  renderer.FillRect(&rect)
  renderer.Present()
}

//doesn't "draw" at screen, but does at the surface to be set at the active surface later
func drawmenuitens(menuitens map[int]string, position int, menumaxitensatscreen int, menu menulist, font *ttf.Font, surface *sdl.Surface) (*sdl.Surface, *sdl.Surface) {
  var (
    buffer bytes.Buffer
    menutext *sdl.Surface
    titletext *sdl.Surface
  )

  //this way we'll end up having as much itens can be at screen starting by "position"
  item := position
  for n := 1; n <= menumaxitensatscreen; n++ {
    //the function that "writes" only works with a single string (as far as i could test it) 
    //we need to set the menu text as it with "\n" as line breaker 

    if menu.Options[menuitens[item]].Check != "" {
      //TODO: Find a better/more secure way to do that (run the command at the yaml)
      if err := exec.Command("/bin/bash", "-c", menu.Options[menuitens[item]].Check).Run(); err != nil {
        buffer.WriteString(menu.Options[menuitens[item]].Desc)
      } else {
        buffer.WriteString(menu.Options[menuitens[item]].Desc2)
      }
    } else {
      buffer.WriteString(menu.Options[menuitens[item]].Desc)
    }
    buffer.WriteString("\n")
    item++
  }
  menuitenstext := buffer.String()

  //this renders our title and menu itens at a surface that will be rendered at the screen later
  titletext, _ = font.RenderUTF8Blended(menu.Name, sdlcolornames.Grey) 
  menutext, _ = font.RenderUTF8BlendedWrapped(menuitenstext, sdlcolornames.Black, int(winWidth - 20))

  return titletext, menutext
}

func runcommand(command string) {
  cmd := exec.Command("/bin/bash", "-c", command)
  err := cmd.Run()
  if err != nil {
    fmt.Println(err)
  }
}

func main() {
  if len(os.Args) != 2 {
    fmt.Printf("Usage: %s <yamlfile>\n", os.Args[0])
  } else {
    os.Exit(run())
  }
}
