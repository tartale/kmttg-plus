import React from "react";

const UpperCenter = ({children}: any) => (
  <div
    style={{
      position: "fixed",
      top: 0,
      padding: "1rem",
    }}
  >
    {children}
  </div>
);

export default UpperCenter;
