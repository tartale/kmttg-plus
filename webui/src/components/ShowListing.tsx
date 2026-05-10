import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableContainer from "@mui/material/TableContainer";
import { useEffect, useState } from "react";
import "./ShowListing.css";
import { mergeEpisodes } from "./showListingHelpers";
import { ShowHeader, ShowRow } from "./ShowRow";
import { Show } from "../services/generated/graphql-types"
import "./TivoStyle.css";
import { useQuery, gql } from '@apollo/client';

const GET_SHOWS = gql`
  query GetShows {
    tivos(filters: [
      {
        name: {eq: "Family Room"}
      }
    ]) {
      name
      shows(filters: [
        {
          title: {
            eq: "60 Minutes"
          }
        }
      ], offset: 0, limit: 25) {
        id
        kind
        title
        description
        recordedOn
        ...on Series {
          episodes {
            id
            seasonNumber
            episodeNumber
            description
            recordedOn
          }
        }
      }
    } 
  }
`;

export type ShowSortField = 'kind' | 'title' | 'recordedOn';

export interface Episode extends Show {
  originalAirDate: Date;
  seasonNumber: number;
  episodeNumber: number;
  episodeTitle: string;
  episodeDescription: string;
}

export interface Movie extends Show {
  movieYear: number;
}

export interface Series extends Show {
  episodes: Episode[];
}

export default function ShowListing(props: any) {
  const [shows, setShows] = useState<Show[]>([]);
  const { data, loading, error } = useQuery(GET_SHOWS);

  useEffect(() => {
    if (data) {
      const tivo = data.tivos.find((t: any) => t.name === "Family Room");
      if (tivo && tivo.shows) {
        const mergedShows = mergeEpisodes(tivo.shows);
        setShows(mergedShows);
      }
    }
  }, [data]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <TableContainer
      component={Paper}
      sx={{ background: "linear-gradient(to bottom, #162c4f, #000000);" }}
      {...props}
    >
      <Table className="showListingTable">
        <ShowHeader/>
        <TableBody>
          {shows.map((show) => (
            <ShowRow key={show.id} show={show} />
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
