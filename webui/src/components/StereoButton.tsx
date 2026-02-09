import React, {useState} from "react";
import "./StereoButton.css";

interface StereoButtonProps {
  label: string;
  [key: string]: any;
}

const StereoButton = ( { label, ...props }: StereoButtonProps) => {
  const [on, setOn] = useState(false);

  const handleClick = () => {
    setOn(!on);
  };

  return (
    <button
      className={`stereo-button ${on ? "on" : "off"}`}
      onClick={handleClick}
      {...props}
    >
      <div className="face">
        <div className="light"></div>
        <div className="indicator"></div>
        <p className="label">{label}</p>
      </div>
    </button>
  );
};

export default StereoButton;
