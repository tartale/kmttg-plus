import React, {useEffect, useState} from "react";
import graphql from "babel-plugin-relay/macro";
import type { TivoSelectorQuery as TivoSelectorQueryType } from "./__generated__/TivoSelectorQuery.graphql";
import { useLazyLoadQuery } from "react-relay";
import StereoButton from "./StereoButton";
import { Box } from "@mui/system";

const TivoSelectorQuery = graphql`
  query TivoSelectorQuery {
    tivos {
      name
    }
  }
`;

function TivoSelector(props: any) {
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

export default TivoSelector;
