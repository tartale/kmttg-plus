import React from "react";

const TiVoLogo = (props: any) => (
    <img {...props} 
      src={`${process.env.PUBLIC_URL}/images/tivo-logo-transparent.png`}
      alt="TiVo logo"
    />
);

export default TiVoLogo;
