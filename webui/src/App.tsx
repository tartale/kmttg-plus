import React, {Suspense} from "react";
import "./App.css";
import "./components/TivoStyle.css";
import ShowListing from "./components/ShowListing";
import TiVoLogo from "./components/TivoLogo";
import TivoSelector from "./components/TivoSelector";
import Loading from "./components/Loading";

function App() {
  return (
    <Suspense fallback={<Loading/>}>
      <TiVoLogo style={{position: "absolute", top: 10, left: 10, width: "100px"}} />
        {/* <TivoSelector style={{position: "absolute", top: 10, left: 20}}/> */}
        <ShowListing style={{position: "absolute", top: 150, left: 10}}/>
    </Suspense>
  );
}

export default App;
