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
      <ShowListing />
    </div>
  );
}

export default App;
