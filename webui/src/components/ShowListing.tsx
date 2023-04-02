import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableContainer from "@mui/material/TableContainer";
import React, { useEffect, useState } from "react";
import { getShows } from "./showListingHelpers";
import "./ShowListing.css";
import "./TivoStyle.css";
import { ShowRow } from "./ShowRow";

export enum ShowKind {
  Movie,
  Series,
  Episode
}

export interface Show {
  recordingId: string;
  kind: ShowKind;
  title: string;
  recordedOn: Date;
  description: string;
}

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

export default function ShowListing() {
  const [shows, setShows] = useState<Show[]>([]);

  useEffect(getShows(setShows), []);

  return (
    <TableContainer
      component={Paper}
      sx={{ background: "linear-gradient(to bottom, #162c4f, #000000);" }}
    >
      <Table className="showListingTable">
        <TableBody>
          {shows.map((show) => (
            <ShowRow key={show.recordingId} show={show} />
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
