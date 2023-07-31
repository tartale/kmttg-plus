import React from "react";
import "./App.css";
import Home from "./components/Home";
import {BrowserRouter, Route, Routes} from 'react-router-dom';

function App() {
  return (
    <BrowserRouter basename="/">
      <Routes>
        <Route path="/" Component={Home}/>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
