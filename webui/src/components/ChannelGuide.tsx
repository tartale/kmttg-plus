import React from "react";
import "./ChannelGuide.css";

function ChannelGuide(props: any) {
  return (
    <div className="table-container">
      <table>
        <thead>
          <tr>
            <th>Channel</th>
            <th>Program</th>
            <th>Time</th>
          </tr>
        </thead>
        <tbody>
          {props.data.map((row: any, index: any) => (
            <tr key={index}>
              <td>{row.channel}</td>
              <td>{row.program}</td>
              <td>{row.time}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default ChannelGuide;
