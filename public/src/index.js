/// app.js
import ReactDOM from "react-dom";
import React from 'react';
import DeckGL from '@deck.gl/react';
import { MVTLayer } from '@deck.gl/geo-layers';
import { StaticMap } from 'react-map-gl';

const interpolate = require('color-interpolate');
const colorParse = require('color-parse')

// Set your mapbox access token here
const MAPBOX_ACCESS_TOKEN = 'pk.eyJ1IjoibWFrZXVwc29tZXRoaW5nIiwiYSI6ImNrNXhqdm5uMjI2b20za29uYnJ6d2NhOHAifQ.RjxFIZyycElU9LqxN7t_MQ';

// Viewport settings
const INITIAL_VIEW_STATE = {
  longitude: 137.97834760987658,
  latitude: 37.55117739361476,
  zoom: 5,
  pitch: 0,
  bearing: 0
};

const prefmap = { 
  "Japan": "117", 
  // "Shiga_Kinki": "4294967814", 
  // "Tochigi_Kant": "4294967817", 
  // "Gunma_Kant": "4294967788", "Yamanashi_Chubu": "4294967825", "Nara_Kinki": "4294967806", "Iwate_Tohoku": "4294967794", "Kagawa_Shikoku": "4294967795", "Kagoshima_Kyushu": "4294967796", "Nagasaki_Kyushu": "4294967805", "Shimane_Chugoku": "4294967815", "Shizuoka_Chubu": "4294967816", "Aichi_Chubu": "4294967779", "Akita_Tohoku": "4294967780", "Aomori_Tohoku": "4294967781", "Kanagawa_Kant": "4294967797", "Kochi_Shikoku": "4294967798", "Kumamoto_Kyushu": "4294967799", "Kyoto_Kinki": "4294967800", "Mie_Kinki": "4294967801", "Miyagi_Tohoku": "4294967802", "Miyazaki_Kyushu": "4294967803", "Nagano_Chubu": "4294967804", "Hiroshima_Chugoku": "4294967789", "Chiba_Kant": "4294967782", "Fukui_Chubu": "4294967784", "Ehime_Shikoku": "4294967783", "Fukuoka_Kyushu": "4294967785", "Gifu_Chubu": "4294967787", "Fukushima_Tohoku": "4294967786", "Tokushima_Shikoku": "4294967818", "Tokyo_Kant": "4294967819", "Tottori_Chugoku": "4294967820", "Toyama_Chubu": "4294967821", "Ibaraki_Kant": "4294967792", "Hyogo_Kinki": "4294967791", 
  // "Hokkaido_Hokkaid": "4294967790", 
  // "Ishikawa_Chubu": "4294967793", "Niigata_Chubu": "4294967807", "Oita_Kyushu": "4294967808", "Okayama_Chugoku": "4294967809", "Osaka_Kinki": "4294967811", "Saga_Kyushu": "4294967812", "Okinawa": "4294967810", "Yamagata_Tohoku": "4294967823", "Wakayama_Kinki": "4294967822", "Yamaguchi_Chugoku": "4294967824", "Saitama_Kant": "4294967813" 
}


const valMin = 0;
const valMax = 4;
// const valMin = 113441.78820385407;
// const valMax = 183819.45159208766;
// const valMin = 0;
// const valMax = 192;
const interpolateFillColor = interpolate(['#225560', "#00AF54", '#EDF060']);
// const host = 'https://dev.tiles.synspective.io/v1/ghi'
const host = 'http://localhost:7778/v1/ghi'
const layers = [];

const colors = []
for (let i=0; i<=valMax; i++) {
  // console.log(d.properties)
  // var v = Math.min((d.properties.intensity-valMin) / (valMax-valMin), 1)
  var c = interpolateFillColor(i/valMax)
  var c2 = colorParse(c)
  // console.log(d.properties.ghi_pred, v)
  colors[i] = c2.values.concat(255)
}

console.log(colors)

for (let i in prefmap) {
  layers.push(new MVTLayer({
    id: `ghi-${prefmap[i]}`,
    data: `${host}/${prefmap[i]}/20180801.z1/{z}/{x}/{y}`,
    // data: `https://dev.tiles.synspective.io/v1/nightlight/4294967819/20180801/{z}/{x}/{y}`,
    // data: `http://localhost:7778/v1/flood-monitoring/4294967296120/20191012/{z}/{x}/{y}`,
    // data: `http://localhost:7778/v1/ghi/4294967821/20190301/{z}/{x}/{y}`,
    // data: `http://localhost:7778/20191012/{z}/{x}/{y}`,
    // pickable: true,
    // stroked: false,
    // filled: true,
    // extruded: true,
    // lineWidthScale: 20,
    // lineWidthMinPixels: 2,
    getFillColor: d => {
      return colors[d.properties.intensity]
    },
    // getLineColor: d => "#f00",
    // getRadius: 400,
    // getLineWidth: 1,
    // getElevation: 30
  }));
}

function App() {
  return (
    <DeckGL
      initialViewState={INITIAL_VIEW_STATE}
      controller={true}
      layers={layers}
    >
      <StaticMap 
          mapStyle="mapbox://styles/makeupsomething/ckd71y9vf0igt1il8f6ll0gu5"
          mapboxApiAccessToken={MAPBOX_ACCESS_TOKEN}
        />
    </DeckGL>
  );
}

ReactDOM.render(<App />, document.getElementById("root"));
