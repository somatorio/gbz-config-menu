package main

import (
  "io/ioutil"
  "fmt"
  "os"
  "time"
  "encoding/xml"
  "sort"
  "bytes"

  "github.com/veandco/go-sdl2/sdl"
  "github.com/veandco/go-sdl2/ttf"
  "gopkg.in/yaml.v2"

  //this makes our life easier by giving the sdl.Color corresponding to the colorname (converted from golang.org/x/image/colornames)
  "sdlcolornames"
)

type (
  items struct {
    Desc string
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
  winWidth, winHeight int32 = 640, 480
  fonttype string = "resources/opensans_hebrew_condensed_regular.ttf"
  fontsize int = 32
  backgroundcolor sdl.Color = sdlcolornames.Whitesmoke
  menucursorcolor sdl.Color = sdlcolornames.Grey
  positiontitle sdl.Rect = sdl.Rect{80, 5, 0, 0}
  positionmenu sdl.Rect = sdl.Rect{30, 55, 0, 0}
  //as the menu cursor height will be set later as the font line height, it doesn't really matter now (so we set as 0 for the lulz =p)
  menucursor sdl.Rect = sdl.Rect{10, positionmenu.Y,620,0}
)
func run() int {
  type (
    Es_input_button struct {
      Name  string `xml:"name,attr"`
      Type  string `xml:"type,attr"`
      Id    uint8 `xml:"id,attr"`
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
      Id    uint8
      Value uint8
    }
  )

  var (
    window *sdl.Window
    font *ttf.Font
    surface *sdl.Surface
    menutext *sdl.Surface
    titletext *sdl.Surface
    input Es_input
    joystick_up, joystick_down, joystick_left, joystick_right, joystick_a, joystick_b Joystick_buttons
    buffer bytes.Buffer
    menuitens map[int]string
  )

  xmlfile, _ := os.Open("es_input.cfg")
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

  //var displaybounds sdl.Rect
  //sdl.GetDisplayBounds(0, &displaybounds)

  window, _ = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
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
  
  //we need a font to work with =p for now we'll use same font/size for title and itens
  font, _ = ttf.OpenFont(fonttype, fontsize)
  
  //so we can set the title to be centralized :)
  titlewidth, _, _ := font.SizeUTF8(menu.Name)
  positiontitle.X = int32(int(winWidth)/2 - titlewidth/2)

  //so we can set the menucursor height to be relative to font size :)
  menucursor.H = int32(font.Height())

  //we need to know how many itens fit at the screen size, it's all math: maxitems = (screenheight-titleheight)/lineheight
  titleheight := int(positiontitle.Y)+int(positionmenu.Y-positiontitle.Y)
  menumaxitensatscreen := (int(winHeight)-titleheight)/font.Height()

  //as the function that "writes" only works with a single string (as far as i could test it) 
  //we need to set the menu text as it with "\n" as line breaker 
  for n := 1; n <= menumaxitensatscreen; n++ {
    buffer.WriteString(menu.Options[menuitens[n]].Desc)
    buffer.WriteString("\n")
  }
  menuitenstext := buffer.String()

  //this renders our title and menu itens at a surface that will be rendered at the screen later
  titletext, _ = font.RenderUTF8Blended(menu.Name, sdlcolornames.Grey) 
  menutext, _ = font.RenderUTF8BlendedWrapped(menuitenstext, sdlcolornames.Black, int(winWidth - 20))
  
  //we need to know the active surface at screen so we can "paste" our texts there later
  surface, _ = window.GetSurface()

  //create the renderer that will colour our background and menu cursor
  renderer, _ := sdl.CreateSoftwareRenderer(surface)

  last := time.Now()
  running := true
  positionmenucursor := 1
  for running {
    renderer.SetDrawColor(backgroundcolor.R, backgroundcolor.G, backgroundcolor.B, backgroundcolor.A)
    renderer.Clear()
    renderer.Present()

    drawmenucursor(renderer, menucursorcolor, menucursor)

    titletext.Blit(nil, surface, &positiontitle)
    menutext.Blit(nil, surface, &positionmenu)

    if time.Since(last).Seconds() > 0.01 {
    event := sdl.PollEvent()
    switch t := event.(type) {
      case *sdl.QuitEvent:
        running = false
      case *sdl.KeyboardEvent:
        if t.Keysym.Sym == sdl.K_q {
          running = false
        }
        if t.Type == sdl.KEYDOWN {
          if t.Keysym.Sym == sdl.K_UP {
            if menucursor.Y != positionmenu.Y {
              menucursor.Y -= menucursor.H
              positionmenucursor--
              last = time.Now()
            }
          }
          if t.Keysym.Sym == sdl.K_DOWN {
            if positionmenucursor < menumaxitensatscreen && positionmenucursor < len(menuitens) {
              menucursor.Y += menucursor.H
              positionmenucursor++
              last = time.Now()
            }
          }
          if t.Keysym.Sym == sdl.K_RETURN {
            fmt.Println(menu.Options[menuitens[positionmenucursor]].Cmd)
          }
        }
      case *sdl.JoyHatEvent:
        if input.Config.Type == "joystick" {
          switch t.Value {
            case joystick_up.Value:
              if menucursor.Y != positionmenu.Y {
                menucursor.Y -= menucursor.H
                last = time.Now()
              }
            case joystick_down.Value:
              menucursor.Y += menucursor.H
              last = time.Now()
          }
        }
      case *sdl.JoyButtonEvent:
        fmt.Printf("%+v\n", t)
      case *sdl.JoyDeviceEvent:
        if t.Type == sdl.JOYDEVICEADDED {
          sdl.JoystickOpen(t.Which)
        }
      }
    }
    window.UpdateSurface()
  }
  return 0
}

func menulistyaml() menulist{
  var menu menulist

  yamlfile, _ := os.Open(os.Args[1])
  menulistfile, _ := ioutil.ReadAll(yamlfile)

  _ = yaml.Unmarshal([]byte(menulistfile),&menu)

  return menu
}

func drawmenucursor(renderer *sdl.Renderer, menucursorcolor sdl.Color, rect sdl.Rect) {
  renderer.SetDrawColor(menucursorcolor.R, menucursorcolor.G, menucursorcolor.B, menucursorcolor.A)
  renderer.FillRect(&rect)
  renderer.Present()
}

func main() {
  if len(os.Args) != 2 {
    fmt.Printf("Usage: %s <yamlfile>\n", os.Args[0])
  } else {
    os.Exit(run())
  }
}
