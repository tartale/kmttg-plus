import TableCell from "@mui/material/TableCell";
import TableRow from "@mui/material/TableRow";
import React from "react";
import { getImageFileForShow, getTitleExtension, recordedOn } from "./showListingHelpers";
import { ShowKind, Series } from "./ShowListing";

export function ShowRow(props: any) {
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

  if (!open ||
    show.kind !== ShowKind.Series) {
    return <React.Fragment />;
  }

  return (
    <React.Fragment>
      {(show as Series).episodes?.map((episode) => (
        <EpisodeRow
          key={episode.recordingId}
          episodeID={episode.recordingId}
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
    <TableCell width={"20%"} style={{ whiteSpace: 'normal' }}>{show.title} {getTitleExtension(show)} </TableCell>
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
      <TableCell width={"10%"}>{recordingDate} {recordingTime}</TableCell>
    </React.Fragment>
  );
}
