import React from "react";
import "./App.css";
import DarkTheme from "./components/DarkTheme";
import TiVoLogo from "./components/TivoLogo";
import TivoSelector from "./components/TivoSelector";
import TivoStyle from "./components/TivoStyle";
import UpperCenter from "./components/UpperCenter";
import UpperLeft from "./components/UpperLeft";

function handleDropdownChange(value: any) {
  console.log("Selected value:", value);
}

function App() {
  return (
    <DarkTheme>
      <div style={TivoStyle}>
        <UpperLeft>
          <TiVoLogo />
        </UpperLeft>
        <UpperCenter>
          <TivoSelector
            file="/config/tivoList.yaml"
            field="tivoList"
            onChange={handleDropdownChange}
          />
        </UpperCenter>
      </div>
    </DarkTheme>
  );
}

export default App;
