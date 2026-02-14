import React, {Suspense} from "react";
import "./TivoStyle.css";
import ShowListing from "./ShowListing";
import TiVoLogo from "./TivoLogo";
import TivoSelector from "./TivoSelector";
import Loading from "./Loading";

function Home() {
  return (
    <Suspense fallback={<Loading/>}>
    <TiVoLogo style={{position: "absolute", top: 10, left: 10, width: "100px"}} />
    <TivoSelector style={{position: "absolute", top: 75, left: 200}}/>
    <ShowListing style={{position: "absolute", top: 150, left: 10}}/>
  </Suspense>
);
}

export default Home;
