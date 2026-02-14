import React, {useEffect, useState} from "react";
import StereoButton from "./StereoButton";
import {
  useQuery,
  gql
 } from "@apollo/client";

const GET_TIVO_NAMES = gql`
 query getTivoNames {
   tivos {
     name
   }
}`;

function TivoSelectorComponent(props: any) {
  const { names } = props
  const [options, setOptions] = useState<string[]>([]);

  const data = useLazyLoadQuery
    <TivoSelectorQueryType>
    (TivoSelectorQuery, {});
  const names: string[] = data.tivos.map((tivo) => tivo.name);

  useEffect(() => {
    setOptions(names);
  }, []);
  
  return (
    <React.Fragment>
      {options.map((option: string, index: number) => (
        <StereoButton label={option} key={index} {...props}/>
      ))}
    </React.Fragment>
  );
}

const TivoSelector = (props: any) => {
  const { loading, error, data } = useQuery(GET_TIVO_NAMES);

  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    console.error(error);
    return <div>Error!</div>;
  }

  const names = data.tivos.map((tivo: { name: string }) => tivo.name);
  return <TivoSelectorComponent names={names} {...props} />;
};

export default TivoSelector;
