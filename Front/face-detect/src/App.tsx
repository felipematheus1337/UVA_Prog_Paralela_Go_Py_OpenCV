import React, { useState, useEffect } from 'react';
import go from "./assets/go.jpg";
import openCV from "./assets/openCV.png";
import py from "./assets/py.jpg";
import './App.css';

function App() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [imageData, setImageData] = useState<string | null>(null);
  const [responseMessage, setResponseMessage] = useState<string | null>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files && event.target.files[0];
    setSelectedFile(file);
  };

  useEffect(() => {
    const sendImage = async () => {
      if (selectedFile) {
        const formData = new FormData();
        formData.append("file", selectedFile);

        try {
          const response = await fetch("http://localhost:8081/upload/single", {
            method: "POST",
            body: formData
          });
          console.log(response)
          if (response.ok) {
            const responseData = await response.text();
            console.log(responseData)
            setResponseMessage(responseData);
            setImageData(URL.createObjectURL(selectedFile));
          } else {
            setResponseMessage("Falha ao enviar arquivo.");
          }
        } catch (error) {
          setResponseMessage("Erro ao enviar arquivo.");
        }
      }
    };

    sendImage();
  }, [selectedFile]);

  const handleRemoveClick = () => {
    setSelectedFile(null);
    setImageData(null);
    setResponseMessage(null);
  };

  return (
    <div className='container'>
      <div className='header'>
        <img src={go} alt="Golang" className="move"/>
        <img src={openCV} alt="OpenCV" className="move" />
        <img src={py} alt="Python" className="move"/>
      </div>

      <div className='content'>
        <input type="file" name="files" onChange={handleFileChange}/>
      </div>

      <section className='result'>
        {responseMessage && <p>{responseMessage}</p>}
        {imageData && (
          <div>
            <img src={imageData} alt="Selected File Preview"/>
            <button className="danger" onClick={handleRemoveClick}>Remover</button>
          </div>
        )}
      </section>
    </div>
  );
}

export default App;
