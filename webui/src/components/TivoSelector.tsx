import React, {useState, useEffect} from "react";
import yaml from "js-yaml";
import TivoStyle from "./TivoStyle";

interface TivoConfig {
  devices: {
    name: string
    address: string
  }[]
}

function TivoSelector(props: any) {
  const [options, setOptions] = useState<string[]>([]);

  useEffect(() => {
    fetch("config/tivo.yaml")
      .then((response) => response.text())
      .then((text) => {
        const data: TivoConfig = yaml.load(text) as TivoConfig;
        const names: string[] = data.devices.map(device => device.name)
        setOptions(names);
      });
  });

  function handleChange(event: any) {
    props.onChange(event.target.value);
  }

  if (!options || options.length === 0) {
    return <div>Loading...</div>;
  }

  return (
    <select onChange={handleChange} style={TivoStyle}>
      {options.map((option, index) => (
        <option key={index} value={option}>
          {option}
        </option>
      ))}
    </select>
  );
}

export default TivoSelector;
