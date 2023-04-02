import React, { useEffect, useMemo, useState } from "react";
// import Paper from "@mui/material/Paper";
// import Table from "@mui/material/Table";
// import TableBody from "@mui/material/TableBody";
// import TableContainer from "@mui/material/TableContainer";
// import TableHead from "@mui/material/TableHead";
// import TableRow from "@mui/material/TableRow";
// import TableCell from "@mui/material/TableCell";
// import TableSortLabel from "@mui/material/TableSortLabel";
// import { useTable, useSortBy, useFilters } from "react-table";
// import { getShows } from "./showListingHelpers";
// import { ShowRow } from "./ShowRow";
// import "./TivoStyle.css";

// export enum ShowKind {
//   Movie,
//   Series,
//   Episode
// }

// export interface Show {
//   recordingId: string;
//   kind: ShowKind;
//   title: string;
//   recordedOn: Date;
//   description: string;
// }

// export interface Movie extends Show {
//   movieYear: number;
// }

// export interface Series extends Show {
//   episodes: Episode[];
// }

// export interface Episode extends Show {
//   originalAirDate: Date;
//   seasonNumber: number;
//   episodeNumber: number;
//   episodeTitle: string;
//   episodeDescription: string;
// }

// export default function ShowListing() {
//   const [shows, setShows] = useState<Show[]>([]);

//   useEffect(() => {
//     getShows(setShows);
//   }, []);

//   const columns = useMemo(
//     () => [
//       {
//         Header: "Title",
//         accessor: "title"
//       },
//       {
//         Header: "Recorded On",
//         accessor: "recordedOn"
//       },
//       {
//         Header: "Description",
//         accessor: "description"
//       }
//     ],
//     []
//   );

//   const tableData = useMemo(() => shows, [shows]);

//   const {
//     getTableProps,
//     getTableBodyProps,
//     headerGroups,
//     rows,
//     prepareRow
//   } = useTable(
//     {
//       columns,
//       data: tableData
//     },
//     useFilters,
//     useSortBy
//   );

//   return (
//     <TableContainer component={Paper}>
//       <Table {...getTableProps()}>
//         <TableHead>
//           {headerGroups.map(headerGroup => (
//             <TableRow {...headerGroup.getHeaderGroupProps()}>
//               {headerGroup.headers.map(column => (
//                 <TableCell
//                   {...column.getHeaderProps(column.getSortByToggleProps())}
//                 >
//                   {column.render("Header")}
//                   <TableSortLabel
//                     active={column.isSorted}
//                     direction={column.isSortedDesc ? "desc" : "asc"}
//                   />
//                 </TableCell>
//               ))}
//             </TableRow>
//           ))}
//         </TableHead>
//         <TableBody {...getTableBodyProps()}>
//           {rows.map(row => {
//             prepareRow(row);
//             return <ShowRow key={row.original.recordingId} show={row.original} />;
//           })}
//         </TableBody>
//       </Table>
//     </TableContainer>
//   );
// }
