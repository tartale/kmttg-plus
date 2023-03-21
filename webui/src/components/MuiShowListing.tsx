import React, { useState, useEffect } from "react";
import Box from '@mui/material/Box';
import Collapse from '@mui/material/Collapse';
import Icon from '@mui/material';
import IconButton from '@mui/material/IconButton';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import { FolderOpenRounded, FolderRounded } from "@mui/icons-material";
import "./TivoStyle.css";
import "./MuiShowListing.css";
import { lightBlue } from "@mui/material/colors";

export interface Show {
  id?: number;
  title: string;
  recordedOn: Date;
  episodes?: Episode[];
  movieYear?: number;
}

export interface Episode {
  id?: number;
  recordedOn?: Date;
  firstAiredOn?: Date;
  season?: number;
  episode?: number;
  episodeTitle?: string;
}

function Row(props: { row: Show }) {
  const { row } = props;
  const [open, setOpen] = React.useState(false);

  const episodeCount = row.episodes?.length || 0
  const episodeCountLabel = episodeCount > 1 ? `[${episodeCount}]` : ""
  const isMovie = row.movieYear || false
  const dayOfWeek = row.recordedOn.toLocaleDateString('en-US', { weekday: 'short' });
  const monthDay = row.recordedOn.toLocaleDateString('en-US', { month: '2-digit', day: '2-digit' });

  const ShowIcon = (props: { episodeCount: number }) => {
    if (episodeCount > 1) {
      return (
        open ?
          <img
            src={"./images/folder-open.png"}
            style={{ width: "3rem", height: "3rem" }}
            alt=""
          />
          :
          <img
            src={"./images/folder-closed.png"}
            style={{ width: "3rem", height: "3rem" }}
            alt=""
          />

      )
    } else {
      const icon = isMovie ? "./images/movie.png" : "./images/television.png"
      return (
        <img
          src={icon}
          style={{ width: "3rem", height: "3rem" }}
          alt=""
        />
      )
    }
  }

  return (
    <React.Fragment>
      <TableRow sx={{ '& > *': { borderBottom: 'unset' } }} onClick={() => setOpen(!open)} >
        <TableCell>
          <ShowIcon episodeCount={episodeCount} />
        </TableCell>
        <TableCell>{row.title} {episodeCountLabel}</TableCell>
        <TableCell>{dayOfWeek}</TableCell>
        <TableCell>{monthDay}</TableCell>
      </TableRow>
      {/* <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box sx={{ margin: 1 }}>
              <Typography variant="h6" gutterBottom component="div">
                History
              </Typography>
              <Table size="small" aria-label="purchases">
                <TableHead>
                  <TableRow>
                    <TableCell>Date</TableCell>
                    <TableCell>Customer</TableCell>
                    <TableCell align="right">Amount</TableCell>
                    <TableCell align="right">Total price ($)</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {row.history.map((historyRow) => (
                    <TableRow key={historyRow.date}>
                      <TableCell component="th" scope="row">
                        {historyRow.date}
                      </TableCell>
                      <TableCell>{historyRow.customerId}</TableCell>
                      <TableCell align="right">{historyRow.amount}</TableCell>
                      <TableCell align="right">
                        {Math.round(historyRow.amount * row.price * 100) / 100}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow> */}
    </React.Fragment>
  );
}

export default function ShowListing() {
  const [shows, setShows] = useState<Show[]>([]);

  useEffect(() => {
    fetch("input/shows.json")
      .then((response) => response.json())
      .then((data) => {
        const parsedShows = data.map((show: Show) => ({
          ...show,
          recordedOn: new Date(show.recordedOn),
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
            <Row key={show.id} row={show} />
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
