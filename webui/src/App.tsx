import React from "react";
import "./App.css";
import Home from "./components/Home";
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import { ApolloClient, ApolloProvider, InMemoryCache, HttpLink } from '@apollo/client';
import { useState } from "react";

const createApolloClient = () => {
 return new ApolloClient({
   link: new HttpLink({
     uri: 'http://localhost:8080/api/query',
   }),
   cache: new InMemoryCache(),
 });
};

function App() {
  const [client] = useState(createApolloClient());
  return (
    <ApolloProvider client={client}>
      <BrowserRouter basename="/">
        <Routes>
          <Route path="/" Component={Home}/>
        </Routes>
      </BrowserRouter>
    </ApolloProvider>
  );
}

export default App;
