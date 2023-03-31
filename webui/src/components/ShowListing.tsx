import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableRow from "@mui/material/TableRow";
import React, { useEffect, useState } from "react";
import { getImageFileForShow, getShows, getTitleExtension, recordedOn } from "./showListingHelpers";
import "./ShowListing.css";
import "./TivoStyle.css";

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

function IconCell(props: any) {
  const { show, open, indent } = props;
  const imageFile: string = getImageFileForShow(show, open);

  const style = indent
    ? { paddingLeft: "2rem", width: "3rem", height: "3rem" }
    : { width: "3rem", height: "3rem" };

  return (
    <TableCell {...props}>
      <img src={imageFile} style={style} alt="" />
    </TableCell>
  );
}

function TitleCell(props: any) {
  const { show } = props;

  return (
    <TableCell width={"20%"} style={{ whiteSpace: 'normal' }}>{show.title} {getTitleExtension(show)} </TableCell>
  )
}

function DescriptionCell(props: any) {
  const { show } = props;

  switch (show.kind) {
    case ShowKind.Movie:
    case ShowKind.Episode:
      return (
        <TableCell {...props} style={{ whiteSpace: 'normal' }} >{show.description}</TableCell>
      )
    default:
      return (<TableCell {...props} style={{ whiteSpace: 'normal' }} />)
  }
}

function RecordingDateCell(props: any) {
  const { show } = props;
  const recordingDate = recordedOn(show)?.toLocaleDateString("en-US", {
    dateStyle: "medium"
  });
  const recordingTime = recordedOn(show)?.toLocaleTimeString("en-US", {
    timeStyle: "short"
  });

  return (
    <React.Fragment>
      <TableCell width={"10%"} >{recordingDate} {recordingTime}</TableCell>
    </React.Fragment>
  )
}

function Row(props: any) {
  const { show } = props;
  const [open, setOpen] = React.useState(false);

  return (
    <React.Fragment>
      <TableRow onClick={() => setOpen(!open)}>
        <IconCell show={show} open={open} indent={false} width={"5%"} />
        <TitleCell show={show} />
        <DescriptionCell show={show} />
        <RecordingDateCell show={show} />
      </TableRow>
      <EpisodeRows show={show} open={open} />
    </React.Fragment>
  );
}

function EpisodeRow(props: any) {
  const { episodeID, show } = props;

  return (
    <TableRow key={episodeID} className="indented">
      <IconCell show={show} open={false} indent={true} width={"5%"} />
      <TableCell>{show.title} {getTitleExtension(show)}</TableCell>
      <DescriptionCell show={show} />
      <RecordingDateCell show={show} />
    </TableRow>
  );
}

function EpisodeRows(props: any) {
  const { show, open } = props;

  if (
    !open ||
    show.kind !== ShowKind.Series
  ) {
    return <React.Fragment />;
  }

  return (
    <React.Fragment>
      {(show as Series).episodes?.map((episode) => (
        <EpisodeRow
          key={episode.recordingId}
          episodeID={episode.recordingId}
          show={{ ...show, ...episode }}
        />
      ))}
    </React.Fragment>
  );
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
            <Row key={show.recordingId} show={show} />
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
