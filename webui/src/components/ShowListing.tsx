import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableContainer from "@mui/material/TableContainer";
import { useEffect, useState } from "react";
import "./ShowListing.css";
import { getShows } from "./showListingHelpers";
import { ShowHeader, ShowRow } from "./ShowRow";
import { Show } from "../services/generated/graphql-types"
import "./TivoStyle.css";

export type ShowSortField = 'kind' | 'title' | 'recordedOn';

export interface Movie extends Show {
  movieYear: number;
}

export interface Series extends Show {
  episodes: Episode[];
}

export interface Episode extends Show {
  originalAirDate: Date;
  seasonNumber: number;
  episodeNumber: number;
  episodeTitle: string;
  episodeDescription: string;
}

export default function ShowListing(props: any) {
  const [shows, setShows] = useState<Show[]>([]);

  useEffect(getShows(setShows), []);

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
            <ShowRow key={show.recordingID} show={show} />
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
