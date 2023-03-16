import React from "react";
import "./App.css";
import DarkTheme from "./components/DarkTheme";
import {Show, ShowListing} from "./components/ShowListing";
import TiVoLogo from "./components/TivoLogo";
import TivoSelector from "./components/TivoSelector";
import TivoStyle from "./components/TivoStyle";
import UpperCenter from "./components/UpperCenter";
import UpperLeft from "./components/UpperLeft";

function handleDropdownChange(value: any) {
  console.log("Selected value:", value);
}

const shows: Show[] = [
  {
    recordedOn: new Date().toISOString(),
    title: "The Big Bang Theory",
    episodeTitle: "The Proposal Proposal"
  },
  {
    recordedOn: new Date().toISOString(),
    title: "Young Sheldon",
    episodeTitle: "A Solar Calculator, a Game Ball, and a Cheerleader's Bosom",
  },
  // ...more shows here
];

function App() {
  return (
    <DarkTheme>
      <div style={TivoStyle}>
        <UpperLeft>
          <TiVoLogo />
        </UpperLeft>
          <TivoSelector
            file="/config/tivoList.yaml"
            field="tivoList"
            onChange={handleDropdownChange}
          />
        <ShowListing shows={shows}/>
      </div>
    </DarkTheme>
  );
}

export default App;
