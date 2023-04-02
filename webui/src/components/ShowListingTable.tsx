import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import { useEffect, useState } from "react";
import { Show } from "./ShowListing";
import "./ShowListing.css";
import { getShows } from "./showListingHelpers";
import { ShowRow } from "./ShowRow";
import "./TivoStyle.css";

export type ShowSortField = 'kind' | 'title' | 'recordedOn';

export default function ShowTable() {
  const [shows, setShows] = useState<Show[]>([]);
  const [sortedShows, setSortedShows] = useState<Show[]>([]);

  useEffect(() => {
    getShows((shows) => {
      setShows(shows);
      setSortedShows(shows);
    });
  }, []);

  const handleSort = (field: ShowSortField) => {
    const sorted = [...sortedShows].sort((a, b) => {
      if (a[field] < b[field]) {
        return -1;
      }
      if (a[field] > b[field]) {
        return 1;
      }
      return 0;
    });
    setSortedShows(sorted);
  };

  const handleFilter = (event: React.ChangeEvent<HTMLInputElement>) => {
    const filtered = shows.filter((show) =>
      show.title.toLowerCase().includes(event.target.value.toLowerCase())
    );
    setSortedShows(filtered);
  };

  return (
    <TableContainer
      component={Paper}
      sx={{ background: "linear-gradient(to bottom, #162c4f, #000000);" }}
    >
      <Table className="showListingTable">
        <TableHead>
          <TableRow>
            <TableCell onClick={() => handleSort("title")}>Title</TableCell>
            <TableCell onClick={() => handleSort("recordedOn")}>Recorded On</TableCell>
            <TableCell>Description</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          <TableRow>
            <TableCell>
              <input type="text" onChange={handleFilter} />
            </TableCell>
          </TableRow>
          {sortedShows.map((show) => (
            <ShowRow key={show.recordingId} show={show} />
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
}
