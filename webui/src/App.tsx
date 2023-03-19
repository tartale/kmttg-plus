import React from "react";
import "./App.css";
import {Show, ShowListing} from "./components/ShowListing";
import TiVoLogo from "./components/TivoLogo";
import TivoSelector from "./components/TivoSelector";

function handleDropdownChange(value: any) {
  console.log("Selected value:", value);
}

const shows: Show[] = [
  {
    recordedOn: new Date().toISOString(),
    title: "The Big Bang Theory",
    episodeTitle: "The Proposal Proposal",
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
    <div>
      <TiVoLogo />
      <TivoSelector onChange={handleDropdownChange} />
      <ShowListing shows={shows} />
    </div>
  );
}

export default App;
