import React, { Component } from 'react';
import * as Pixi from 'pixi.js'
import './App.css';

class App extends Component {
  constructor(props) {
    super(props)

    // TODO: figure out where it's correct to put this in JS
    this.BODY_COLORS = {
      "Carl": 0xFF0000,
      "Frog": 0x00FF00
    }
    this.DEFAULT_BODY_COLOR = 0x0000FF

    // Set up pixi app and and canvas
    this.pixi_cnt = null
    this.app = new Pixi.Application({width: window.innerWidth, height: window.innerHeight, antialias: true, backgroundColor: 0x282c34})
    this.app.renderer.autoResize = true

    // Init background
    window.addEventListener('resize', this.resize)

    // World Transform
    this.world = new Pixi.Container()
    this.world.position.x = window.innerWidth / 2
    this.world.position.y = window.innerHeight / 2
    this.world.scale.y = -1
    this.app.stage.addChild(this.world)

    // Static Elements
    this.static = new Pixi.Container()
    this.world.addChild(this.static)

    // Dynamic Elements
    this.dynamic = new Pixi.Container()
    this.world.addChild(this.dynamic)

    this.bodies = new Pixi.Container()
    this.body_refs = {}
    this.dynamic.addChild(this.bodies)

    this.legs = new Pixi.Container()
    this.dynamic.addChild(this.legs)

    // Init mouse follower
    this.circle = new Pixi.Graphics()
    this.circle.beginFill(0x8b929e)
    this.circle.drawCircle(0, 0, 5)
    this.circle.endFill()
    this.circle.position.set(window.innerWidth / 2, window.innerHeight / 2)
    this.app.stage.addChild(this.circle)

    // Connect to game server
    let target = "ws://" + window.location.hostname + ":8000/connect"
    console.log("App.constructor: connecting on " + target)
    this.conn = new WebSocket(target)
    this.conn.onmessage = this.onMessage

    // State
    this.state = {
      "GameState": {
        "tid": "GameState",
        "bodies": [],
        "legs": [],
        "owner_id": 0
      },
      "MapState": {
        "tid": "MapState",
        "surfaces": []
      }
    }
  }

  resize = () => {
    this.app.renderer.resize(window.innerWidth, window.innerHeight);
    this.world.x = window.innerWidth / 2
    this.world.y = window.innerHeight / 2
  }

  onMessage = (msg) => {
    let data = JSON.parse(msg.data)
    if (!data.tid) {
      console.log("App.onMessage: message did not have tid in response")
    }

    var intermediate_state = this.state
    intermediate_state[data.tid] = data
    this.setState(intermediate_state)

    if (data.tid === "MapState") {
      this.onMapUpdate(data)
    } else if (data.tid === "GameState") {
      this.onGameUpdate(data)
    }
  }

  onMapUpdate = (map) => {
    // Clean up old children
    for (var child in this.static.children) {
      this.static.removeChild(child)
    }

    // TODO: remove
    // Add marking circles
    let xs = [0, -50, -50, 50, 50]
    let ys = [0, -50, 50, -50, 50]
    // white, green, red, yellow, purple
    let clr = [0xFFFFFF, 0x3bdd25, 0xd6135a, 0xe8e009, 0xb412e5]
    for (var i = 0; i < 5; i++) {
      var circ = new Pixi.Graphics()
      circ.beginFill(clr[i]);
      circ.drawCircle(0, 0, 5);
      circ.endFill();
      circ.position.set(xs[i], ys[i])
      this.static.addChild(circ)
    }

    // Add all the map elements
    for (var idx in map.surfaces) {
      var surface = map.surfaces[idx]
      var graphics = new Pixi.Graphics()
      graphics.beginFill(0xAA8844)

      let p3x = surface.p0.x + (surface.p2.x - surface.p1.x)
      let p3y = surface.p0.y + (surface.p2.y - surface.p1.y)
      graphics.drawPolygon([surface.p0.x, surface.p0.y,
                            surface.p1.x, surface.p1.y,
                            surface.p2.x, surface.p2.y,
                            p3x, p3y])
      graphics.endFill()
      this.world.addChild(graphics)
    }
  }

  onGameUpdate = (game) => {
    // Draw bodies first
    // TODO: have bodies be a sprite
    for (var idx in game.bodies) {
      if (game.bodies[idx].name in this.body_refs) {
        this.body_refs[game.bodies[idx].name].position.x = game.bodies[idx].pose.x
        this.body_refs[game.bodies[idx].name].position.y = game.bodies[idx].pose.y
        this.body_refs[game.bodies[idx].name].rotation = game.bodies[idx].pose.theta
      } else {
        var body = new Pixi.Graphics();
        var color = this.DEFAULT_BODY_COLOR
        if (game.bodies[idx].name in this.BODY_COLORS) {
          color = this.BODY_COLORS[game.bodies[idx].name]
        }
        body.beginFill(color)
        body.drawRoundedRect(-30, -10, 60, 20, 10)
        body.endFill();
        body.position.set(game.bodies[idx].pose.x, game.bodies[idx].pose.x)
        body.rotation = game.bodies[idx].pose.theta
        this.bodies.addChild(body)
        this.body_refs[game.bodies[idx].name] = body
      }
    }

    // TODO: draw legs
  }

  onMouseMove = (e) => {
    this.circle.position.set(e.clientX, e.clientY)
  }

  onMouseUp = (e) => {
    console.log("clicked")
  }

  updatePixiCnt = (element) => {
    this.pixi_cnt = element

    if(this.pixi_cnt && this.pixi_cnt.children.length <= 0) {
       this.pixi_cnt.appendChild(this.app.view)
       this.setup()
     }
  }

  setup = () => {
    console.log("TODO: run setup")
  }

  render() {
    return <div onMouseUp={this.onMouseUp} onMouseMove={this.onMouseMove} ref={this.updatePixiCnt} />
  }
}

export default App;
