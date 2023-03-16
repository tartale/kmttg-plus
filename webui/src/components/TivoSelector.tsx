import React, {useState, useEffect} from "react";
import yaml from "js-yaml";
import TivoStyle from "./TivoStyle";

interface TivoList {
  [key: string]: string[];
}

function TivoSelector(props: any) {
  const [options, setOptions] = useState<string[]>([]);

  useEffect(() => {
    fetch(props.file)
      .then((response) => response.text())
      .then((text) => {
        const data: TivoList = yaml.load(text) as TivoList;
        setOptions(data[props.field]);
      });
  }, [props.file, props.field]);

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
