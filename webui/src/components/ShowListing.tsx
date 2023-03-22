import React, {useState, useEffect} from "react";
import Box from "@mui/material/Box";
import Collapse from "@mui/material/Collapse";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableRow from "@mui/material/TableRow";
import Paper from "@mui/material/Paper";
import "./TivoStyle.css";
import "./ShowListing.css";
import {v4 as uuidv4} from "uuid";

export interface Show {
  id: string;
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
  firstAiredOn?: Date;
  season?: number;
  episode?: number;
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

function IconCell(props: {show: Show; open: boolean}) {
  const {show, open} = props;
  const imageFile: string = getImageFileForShow(show, open);

  return (
    <TableCell>
      <img src={imageFile} style={{width: "3rem", height: "3rem"}} alt="" />
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
      <TableRow
        sx={{"& > *": {borderBottom: "unset"}}}
        onClick={() => setOpen(!open)}
      >
        <IconCell show={show} open={open} />
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

function EpisodeRow(props: {key: string; show: Show}) {
  const {key, show} = props;
  const {dayOfWeek, monthDay} = parseRecordingDate(show);

  return (
    <TableRow key={key} sx={{"& > *": {borderBottom: "unset"}}}>
      <IconCell show={show} open={false} />
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
      <TableRow>
        <Collapse in={open} timeout="auto" unmountOnExit>
          <Box sx={{margin: 1}}>
            <Table className="showListingTable">
              <TableBody>
                {(show as Series).episodes?.map((episode) => (
                  <EpisodeRow key={episode.id} show={{...show, ...episode}} />
                ))}
              </TableBody>
            </Table>
          </Box>
        </Collapse>
      </TableRow>
    </React.Fragment>
  );
}

export default function ShowListing() {
  const [shows, setShows] = useState<Show[]>([]);

  useEffect(() => {
    fetch("input/shows.json")
      .then((response) => response.json())
      .then((jsonArray) => {
        const parsedShows = jsonArray.map((obj: any): Show[] => ({
          ...obj,
          kind: (obj as Series).episodes ? "series" : "movie",
          id: obj.id || uuidv4(),
          recordedOn: new Date(obj.recordedOn),
          episodes: obj.episodes?.map(
            (episode: Episode): Episode => ({
              ...episode,
              kind: "episode",
              id: obj.id || uuidv4(),
              recordedOn: new Date(episode.recordedOn),
            })
          ),
        }));
        setShows(parsedShows);
      })
      .catch((error) => console.error(error));
  }, []);

  return (
    <TableContainer component={Paper}>
      <Table className="showListingTable">
        <TableBody>
          {shows.map((show) => (
            <Row key={show.id} show={show} />
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

// {/* <TableRow>
//   <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
//     <Collapse in={open} timeout="auto" unmountOnExit>
//       <Box sx={{ margin: 1 }}>
//         <Typography variant="h6" gutterBottom component="div">
//           History
//         </Typography>
//         <Table size="small" aria-label="purchases">
//           <TableHead>
//             <TableRow>
//               <TableCell>Date</TableCell>
//               <TableCell>Customer</TableCell>
//               <TableCell align="right">Amount</TableCell>
//               <TableCell align="right">Total price ($)</TableCell>
//             </TableRow>
//           </TableHead>
//           <TableBody>
//             {row.history.map((historyRow) => (
//               <TableRow key={historyRow.date}>
//                 <TableCell component="th" scope="row">
//                   {historyRow.date}
//                 </TableCell>
//                 <TableCell>{historyRow.customerId}</TableCell>
//                 <TableCell align="right">{historyRow.amount}</TableCell>
//                 <TableCell align="right">
//                   {Math.round(historyRow.amount * row.price * 100) / 100}
//                 </TableCell>
//               </TableRow>
//             ))}
//           </TableBody>
//         </Table>
//       </Box>
//     </Collapse>
//   </TableCell>
// </TableRow> */}
