import React from "react";

const LightTheme = ({children}: any) => (
  <div
    style={{
      backgroundColor: "#fff",
      color: "#292929",
      padding: "1rem",
      display: "flex",
      flexDirection: "column",
      alignItems: "center",
      justifyContent: "center",
      minHeight: "100vh",
    }}
  >
    {children}
  </div>
);

export default LightTheme;
