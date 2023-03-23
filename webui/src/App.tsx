import React from "react";
import "./App.css";
import "./components/TivoStyle.css";
import ShowListing from "./components/ShowListing";
import TiVoLogo from "./components/TivoLogo";

function App() {
  return (
    <div>
      <div style={{position: "absolute", top: 10, left: 10}}>
        <TiVoLogo />
      </div>
      <div style={{position: "absolute", top: 150, left: 10}}>
        <ShowListing />
      </div>
    </div>
  );
}

export default App;
