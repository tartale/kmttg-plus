import React, {useState} from "react";
import IconButton from "@mui/material/IconButton";
import { graphql } from 'relay-runtime';
import type { TivoSelectorQuery as TivoSelectorQueryType } from "./__generated__/TivoSelectorQuery.graphql";
import { useLazyLoadQuery } from "react-relay";

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
  setOptions(names);
  
  function handleChange(event: any) {
    props.onChange(event.target.value);
  }

  return (
    <select onChange={handleChange}>
       {options.map((option, index) => (
         <option key={index} value={option}>
           {option}
         </option>
       ))}
     </select>
  );
}

export default TivoSelector;
