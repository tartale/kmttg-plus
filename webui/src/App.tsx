import React from "react";
import "./App.css";
import "./components/TivoStyle.css";
import ShowListing from "./components/MuiShowListing"
import TiVoLogo from "./components/TivoLogo";
import TivoSelector from "./components/TivoSelector";

function handleDropdownChange(value: any) {
  console.log("Selected value:", value);
}

function App() {

  return (
    <div>
      <div style={{ position: "absolute", top: 10, left: 10 }}>
        <TiVoLogo />
      </div>
      <div style={{ position: "absolute", top: 100, right: 10}}>
        <TivoSelector onChange={handleDropdownChange} />
      </div>
      <ShowListing />
    </div>
  );
}

export default App;
