import React from "react";

const DarkTheme = ({children}: any) => (
  <div
    style={{
      backgroundColor: "#292929",
      color: "#fff",
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
export default DarkTheme;
