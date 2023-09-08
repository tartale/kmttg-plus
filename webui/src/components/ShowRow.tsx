import TableCell from "@mui/material/TableCell";
import TableRow from "@mui/material/TableRow";
import TableHead from "@mui/material/TableHead";
import React from "react";
import { getImageFileForShow, getTitleExtension } from "./showListingHelpers";
import { ShowKind, Series } from "../services/generated/graphql-types"
import IconButton from "@mui/material/IconButton";
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import PlayCircleIcon from '@mui/icons-material/PlayCircle';
import DownloadSharpIcon from '@mui/icons-material/DownloadSharp';
import { makeStyles } from '@mui/styles';

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
  const indent: boolean = (show.kind === ShowKind.Episode)

  return (
    <React.Fragment>
      <TableRow key={show.recordingID} onClick={() => setOpen(!open)}>
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
          key={episode.id}
          show={{ ...show, ...episode }} />
      ))}
    </React.Fragment>
  );
}

function IconCell(props: any) {
  const { show, open, indent, ...remainingProps } = props;
  const imageFile: string = getImageFileForShow(show, open);

  const style = indent
    ? { paddingLeft: "2rem", width: "3rem" }
    : { width: "3rem" };

  return (
    <TableCell {...remainingProps}>
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
    case ShowKind.Series:
      return (
        <TableCell {...props} style={{ whiteSpace: 'normal' }}>{show.description}</TableCell>
      );
    case ShowKind.Episode:
      return (
        <TableCell {...props} style={{ whiteSpace: 'normal', paddingLeft: "3rem", paddingRight: "10rem" }}>{show.episodeDescription}</TableCell>
      );
    default:
      return (<TableCell {...props} style={{ whiteSpace: 'normal' }} />);
  }
}

function RecordedOnCell(props: any) {
  const { show } = props;
  const recordedOn = new Date(show.recordedOn);
  const recordingDate = recordedOn.toLocaleDateString("en-US", {
    dateStyle: "medium"
  });
  const recordingTime = recordedOn.toLocaleTimeString("en-US", {
    timeStyle: "short"
  });

  return (
    <React.Fragment>
      <TableCell {...props}>{recordingDate} {recordingTime}</TableCell>
    </React.Fragment>
  );
}

const actionCellStyle = makeStyles(() => ({
  iconButton: {
    backgroundColor: 'white',
    margin: '0.25rem',
  },
}));

function ActionCell(props: any) {
  const { show } = props;
  const classes = actionCellStyle();

  switch (show.kind) {
    case ShowKind.Movie:
    case ShowKind.Episode:
      return (
        <React.Fragment>
          <TableCell {...props}>
          <IconButton size="small" className={classes.iconButton}>
              <DownloadSharpIcon/>
            </IconButton>
            <IconButton size="small" className={classes.iconButton}>
              <PlayCircleIcon/>
            </IconButton>
          </TableCell>
        </React.Fragment>
      )
    default:
      return (
        <React.Fragment>
          <TableCell {...props}>
            <IconButton size="medium" className={classes.iconButton}>
              <ExpandMoreIcon/>
            </IconButton>
          </TableCell>
        </React.Fragment>
      )
  }

}
