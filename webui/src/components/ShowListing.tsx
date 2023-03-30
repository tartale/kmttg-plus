import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableRow from "@mui/material/TableRow";
import React, {useEffect, useState} from "react";
import {v4 as uuidv4} from "uuid";
import "./ShowListing.css";
import "./TivoStyle.css";

export interface Show {
  recordingId: string;
  kind: string;
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
  originalAirDate?: Date;
  seasonNumber?: number;
  episodeNumber?: number;
  episodeTitle?: string;
  episodeDescription?: string;
}

export const recordedOn = (show: Show): Date | undefined =>
  (show as Movie).recordedOn ||
  (show as Episode).recordedOn ||
  (show as Series).episodes.reduce(
    (latest, episode) =>
      episode.recordedOn > latest ? episode.recordedOn : latest,
    new Date(0)
  ) ||
  undefined;

export const getImageFileForShow = (show: Show, open: boolean): string => {
  switch (show.kind) {
    case "series": {
      const episodeCount = (show as Series).episodes.length;
      if (episodeCount === 1) {
        return "./images/television.png";
      }
      if (open) {
        return "./images/folder-open.png";
      } else {
        return "./images/folder-closed.png";
      }
    }
    case "episode": {
      return "./images/television.png";
    }
    case "movie": {
      return "./images/movie.png";
    }
    default: {
      return "./images/television-unknown.png";
    }
  }
};

type ShowSetter = React.Dispatch<React.SetStateAction<Show[]>>

const getShows = (setShows: ShowSetter) => () => {
  fetch("http://localhost:8181/getMyShows?limit=50&tivo=Living%20Room&offset=0", {
    "credentials": "omit",
    "headers": {
        "Accept": "application/json, text/javascript, */*; q=0.01",
        "Accept-Language": "en-US,en;q=0.5",
        "X-Requested-With": "XMLHttpRequest",
        "Sec-Fetch-Dest": "empty",
        "Sec-Fetch-Mode": "cors",
        "Sec-Fetch-Site": "same-origin",
        "Sec-GPC": "1",
        "Pragma": "no-cache",
        "Cache-Control": "no-cache",
    },
    "method": "GET",
    "mode": "cors"
  })
    .then((response) => response.json())
    .then((jsonArray) => {
      const shows = mergeEpisodes(jsonArray.map(parseShow))
      setShows(shows);
    })
    .catch((error) => console.error(error));

}

const parseShow = (obj: any): Show => {
  const recording = obj.recording[0]
  const show: Series | Movie = {
    recordingId: recording.recordingId,
    kind: recording.episodic ? "series" : "movie",
    title: recording.title,
    recordedOn: new Date(recording.startTime),
    description: recording.descrtipion,
    movieYear: recording.movieYear,
    episodes: recording.episodic? [
      {
        recordingId: recording.recordingId,
        kind: "episode",
        title: recording.title,
        recordedOn: new Date(recording.startTime),
        description: recording.descrtipion,
        episodeNumber: recording.episodeNum? recording.episodeNum[0] : 0,
        episodeTitle: recording.subtitle,
        episodeDescription: recording.description,
      }
    ]: []
  }

  return show
}

const mergeEpisodes = (shows: Show[]): Show[] => {
  const combinedShows: Show[] = Object.values(shows.reduce((acc: any, show) => {
    if (show.kind === "series") {
      const series = (show as Series)
      if (acc[series.title]) {
        acc[series.title].episodes.push(...series.episodes);
      } else {
        acc[series.title] = {...series};
      }
    } else {
      acc[show.title] = {...show}
    }
    return acc;
  }, {}));
  
  return combinedShows
}

const parseRecordingDate = (show: Show) => {
  const dayOfWeek = recordedOn(show)?.toLocaleDateString("en-US", {
    weekday: "short",
  });
  const monthDay = recordedOn(show)?.toLocaleDateString("en-US", {
    month: "2-digit",
    day: "2-digit",
  });

  return {dayOfWeek, monthDay};
};

function IconCell(props: {show: Show; open: boolean; indent: boolean}) {
  const {show, open, indent} = props;
  const imageFile: string = getImageFileForShow(show, open);

  const style = indent
    ? {paddingLeft: "2rem", width: "3rem", height: "3rem"}
    : {width: "3rem", height: "3rem"};

  return (
    <TableCell>
      <img src={imageFile} style={style} alt="" />
    </TableCell>
  );
}

function Row(props: {show: Show}) {
  const {show} = props;
  const [open, setOpen] = React.useState(false);
  const episodeCount = (show as Series).episodes?.length || 0;
  const episodeCountLabel = episodeCount > 1 ? `[${episodeCount}]` : "";
  const {dayOfWeek, monthDay} = parseRecordingDate(show);

  return (
    <React.Fragment>
      <TableRow onClick={() => setOpen(!open)}>
        <IconCell show={show} open={open} indent={false} />
        <TableCell>
          {show.title} {episodeCountLabel}
        </TableCell>
        <TableCell>{dayOfWeek}</TableCell>
        <TableCell>{monthDay}</TableCell>
      </TableRow>
      <EpisodeRows show={show} open={open} />
    </React.Fragment>
  );
}

function EpisodeRow(props: {episodeID: string; show: Show}) {
  const {episodeID, show} = props;
  const {dayOfWeek, monthDay} = parseRecordingDate(show);

  return (
    <TableRow key={episodeID} className="indented">
      <IconCell show={show} open={false} indent={true} />
      <TableCell>{show.title}</TableCell>
      <TableCell>{dayOfWeek}</TableCell>
      <TableCell>{monthDay}</TableCell>
    </TableRow>
  );
}

function EpisodeRows(props: {show: Show; open: boolean}) {
  const {show, open} = props;

  if (
    !open ||
    show.kind !== "series" ||
    (show as Series).episodes.length <= 1
  ) {
    return <React.Fragment />;
  }

  return (
    <React.Fragment>
      {(show as Series).episodes?.map((episode) => (
        <EpisodeRow
          key={episode.recordingId}
          episodeID={episode.recordingId}
          show={{...show, ...episode}}
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
      sx={{background: "linear-gradient(to bottom, #162c4f, #000000);"}}
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
