import React, { Component } from 'react';
import * as Pixi from 'pixi.js'
import './App.css';

class App extends Component {
  constructor(props) {
    super(props);
    this.pixi_cnt = null;
    this.app = new Pixi.Application({width: window.innerWidth, height: window.innerHeight, antialias: true, backgroundColor: 0x282c34});
    let text = new Pixi.Text('Imagine a spider...',{fontFamily : 'Helvetica', fontSize: 96, fill : 0xAAAAAA, align: "center"});
    text.y = window.innerHeight / 2 - text.height / 2;
    text.x = window.innerWidth / 2 - text.width / 2;
    this.app.stage.addChild(text);
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
    return <div ref={this.updatePixiCnt} />;
  }
}

export default App;
