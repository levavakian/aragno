import React, { Component } from 'react';
import * as Pixi from 'pixi.js'
import * as io from 'socket.io-client'
import './App.css';

class App extends Component {
  constructor(props) {
    super(props);
    this.pixi_cnt = null;
    this.app = new Pixi.Application({width: window.innerWidth, height: window.innerHeight, antialias: true, backgroundColor: 0x282c34});
    this.app.renderer.autoResize = true;

    this.text = new Pixi.Text('Imagine a spider...',{fontFamily : 'Helvetica', fontSize: 96, fill : 0xAAAAAA, align: "center"});
    this.text.y = window.innerHeight / 2 - this.text.height / 2;
    this.text.x = window.innerWidth / 2 - this.text.width / 2;
    this.app.stage.addChild(this.text);

    this.circle = new Pixi.Graphics();
    this.circle.beginFill(0x8b929e);
    this.circle.drawCircle(0, 0, 5);
    this.circle.endFill();
    this.circle.position.set(window.innerWidth / 2, window.innerHeight / 2)
    this.app.stage.addChild(this.circle);

    window.addEventListener('resize', this.resize);
  }

  resize = () => {
    this.app.renderer.resize(window.innerWidth, window.innerHeight);
    this.text.y = window.innerHeight / 2 - this.text.height / 2;
    this.text.x = window.innerWidth / 2 - this.text.width / 2;
  }

  onMouseMove = (e) => {
    this.circle.position.set(e.clientX, e.clientY)
  }

  onMouseUp = (e) => {
    console.log("clicked")
  }

  updatePixiCnt = (element) => {
    this.pixi_cnt = element;

    if(this.pixi_cnt && this.pixi_cnt.children.length<=0) {
       this.pixi_cnt.appendChild(this.app.view);
       this.setup();
     }
  };

  setup = () => {
    console.log("TODO: run setup")
  }

  render() {
    return <div onMouseUp={this.onMouseUp} onMouseMove={this.onMouseMove} ref={this.updatePixiCnt} />;
  }
}

export default App;
