import TableCell from "@mui/material/TableCell";
import TableRow from "@mui/material/TableRow";
import TableHead from "@mui/material/TableHead";
import React from "react";
import { getImageFileForShow, getTitleExtension, recordedOn } from "./showListingHelpers";
import { ShowKind, Series } from "./ShowListing";
import IconButton from "@mui/material/IconButton";
import UnfoldMoreSharpIcon from '@mui/icons-material/UnfoldMoreSharp';
import DownloadSharpIcon from '@mui/icons-material/DownloadSharp';
import { makeStyles } from '@mui/material/styles';

export function ShowHeader(props: any) {

  return (
    <TableHead>
      <TableRow>
        <TableCell width={"5%"} ></TableCell>
        <TableCell width={"20%"} >Title</TableCell>
        <TableCell>Description</TableCell>
        <TableCell width={"10%"} >Recorded On</TableCell>
        <TableCell width={"10%"} ></TableCell>
      </TableRow>
    </TableHead>
  )
}

export function ShowRow(props: any) {
  const { show } = props;
  const [open, setOpen] = React.useState(false);
  const indent = (show.kind === ShowKind.Episode)

  return (
    <React.Fragment>
      <TableRow key={show.recordingId} onClick={() => setOpen(!open)}>
        <IconCell width={"5%"} show={show} open={open} indent={indent} />
        <TitleCell width={"20%"} show={show} />
        <DescriptionCell show={show} />
        <RecordedOnCell width={"10%"} show={show} />
        <ActionCell width={"10%"} show={show} />
      </TableRow>
      <EpisodeRows show={show} open={open} />
    </React.Fragment>
  );
}

function EpisodeRows(props: any) {
  const { show, open } = props;

  if (!open ||
    show.kind !== ShowKind.Series) {
    return <React.Fragment />;
  }

  return (
    <React.Fragment>
      {(show as Series).episodes?.map((episode) => (
        <ShowRow
          key={episode.recordingId}
          show={{ ...show, ...episode }} />
      ))}
    </React.Fragment>
  );
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
    <TableCell {...props} style={{ whiteSpace: 'normal' }}>{show.title} {getTitleExtension(show)} </TableCell>
  );
}

function DescriptionCell(props: any) {
  const { show } = props;

  switch (show.kind) {
    case ShowKind.Movie:
    case ShowKind.Episode:
      return (
        <TableCell {...props} style={{ whiteSpace: 'normal' }}>{show.description}</TableCell>
      );
    default:
      return (<TableCell {...props} style={{ whiteSpace: 'normal' }} />);
  }
}

function RecordedOnCell(props: any) {
  const { show } = props;
  const recordingDate = recordedOn(show)?.toLocaleDateString("en-US", {
    dateStyle: "medium"
  });
  const recordingTime = recordedOn(show)?.toLocaleTimeString("en-US", {
    timeStyle: "short"
  });

  return (
    <React.Fragment>
      <TableCell {...props}>{recordingDate} {recordingTime}</TableCell>
    </React.Fragment>
  );
}

const actionCellStyle = makeStyles((theme: any) => ({
  iconButton: {
    backgroundColor: '#ffffff', // set the background color of the icon button
    borderRadius: '50%', // make the icon button circular
    padding: theme.spacing(2), // increase the padding of the icon button
  },
}));

function ActionCell(props: any) {
  const { show } = props;
  const classes = useStyles();

  var icon
  switch (show.kind) {
    case ShowKind.Movie:
    case ShowKind.Episode:
      icon = (<DownloadSharpIcon/>)
      break
    default:
      icon = (<UnfoldMoreSharpIcon/>)
  }

  return (
    <React.Fragment>
      <TableCell {...props}>
        <IconButton size="large" className={classes.iconButton}>
          {icon}
        </IconButton>
      </TableCell>
    </React.Fragment>
  );
}
