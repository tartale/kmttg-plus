import React, {useEffect, useState, Fragment} from "react";
import StereoButton from "./StereoButton";
import {
  ApolloClient,
  InMemoryCache,
  useQuery,
  gql
 } from "@apollo/client";
 import { Tivo } from "../services/generated/graphql-types"

// const client = new ApolloClient({
//   uri: 'http://localhost:8080/api/query',
//   cache: new InMemoryCache()
//  });

//  async function getNames(): Promise<string[]> {
//   const result = await client.query({
//     query: gql`
//       {
//         tivos {
//           name
//         }
//       }
//     `,
//   });
//   const data: Tivo[] = result.data.tivos;
//   return data.map((tivo) => tivo.name);
// }

const GET_TIVO_NAMES = gql`
 query getTivoNames {
   tivos {
     name
   }
}`;

function TivoSelectorComponent(props: any) {
  const { names } = props
  const [options, setOptions] = useState<string[]>([]);

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
  return <TivoSelectorComponent {...props} names={names} />;
};

export default TivoSelector;
