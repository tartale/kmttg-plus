import React from "react";

const UpperLeft = ({children}: any) => (
  <div
    style={{
      position: "fixed",
      top: 0,
      left: 0,
      padding: "1rem",
    }}
  >
    {children}
  </div>
);
export default UpperLeft;
