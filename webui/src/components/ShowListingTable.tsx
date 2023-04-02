// import Paper from "@mui/material/Paper";
// import Table from "@mui/material/Table";
// import TableBody from "@mui/material/TableBody";
// import TableContainer from "@mui/material/TableContainer";
// import TableRow from "@mui/material/TableRow";
// import TableCell from "@mui/material/TableCell";
// import TableHead from "@mui/material/TableHead";
import { useEffect, useState } from "react";
// import "./ShowListing.css";
// import { getShows } from "./showListingHelpers";
// import { Show, ShowKind } from "./ShowListing";
// import "./TivoStyle.css";

// export default function ShowListing() {
//   const [shows, setShows] = useState<Show[]>([]);
//   const [sortedShows, setSortedShows] = useState<Show[]>([]);

//   useEffect(() => {
//     getShows((shows) => {
//       setShows(shows);
//       setSortedShows(shows);
//     });
//   }, []);

//   const handleSort = (field: string) => {
//     const sorted = [...sortedShows].sort((a, b) => {
//       if (a[field] < b[field]) {
//         return -1;
//       }
//       if (a[field] > b[field]) {
//         return 1;
//       }
//       return 0;
//     });
//     setSortedShows(sorted);
//   };

//   const handleFilter = (event: React.ChangeEvent<HTMLInputElement>) => {
//     const filtered = shows.filter((show) =>
//       show.title.toLowerCase().includes(event.target.value.toLowerCase())
//     );
//     setSortedShows(filtered);
//   };

//   return (
//     <TableContainer
//       component={Paper}
//       sx={{ background: "linear-gradient(to bottom, #162c4f, #000000);" }}
//     >
//       <Table className="showListingTable">
//         <TableHead>
//           <TableRow>
//             <TableCell onClick={() => handleSort("title")}>Title</TableCell>
//             <TableCell onClick={() => handleSort("recordedOn")}>Recorded On</TableCell>
//             <TableCell>Description</TableCell>
//           </TableRow>
//         </TableHead>
//         <TableBody>
//           <TableRow>
//             <TableCell>
//               <input type="text" onChange={handleFilter} />
//             </TableCell>
//           </TableRow>
//           {sortedShows.map((show) => (
//             <TableRow key={show.recordingId}>
//               <TableCell>{show.title}</TableCell>
//               <TableCell>{show.recordedOn.toLocaleDateString()}</TableCell>
//               <TableCell>{show.description}</TableCell>
//             </TableRow>
//           ))}
//         </TableBody>
//       </Table>
//     </TableContainer>
//   );
// }
