import React, { useState } from 'react';
import axios from 'axios';

function AdminPage() {
  const [formData, setFormData] = useState({
    title: '',
    cover_image: '',
    type: '',
    publisher: '',
    mal_id: '',
    score: '',
    popularity: ''
  });

  const [selectedFile, setSelectedFile] = useState(null);

  const handleFileChange = (e) => {
    setSelectedFile(e.target.files[0]);
  };

  const handleUploadAndAddData = async () => {
    // Upload the file
    const formDataFile = new FormData();
    formDataFile.append('file', selectedFile);
    try {
      const response = await axios.post('http://localhost:8080/upload', formDataFile, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      });
      console.log(response.data);
      // Handle successful upload
    } catch (error) {
      console.error('Error:', error);
      // Handle upload error
    }


    const fileName = selectedFile.name;
    const coverImageUrl = `http://localhost:8080/covers/${fileName}`;
    // Add data to the database
    const formDataDB = {
      ...formData,
      cover_image: coverImageUrl // Assuming the server responds with the URL of the uploaded image
    };
    try {
      const addDataResponse = await axios.post('http://localhost:8080/add-data', formDataDB);
      console.log(addDataResponse.data);
      // Provide feedback to the user indicating success
    } catch (error) {
      console.error('Error adding data:', error);
      // Provide feedback to the user indicating failure
    }
  };

  const handleChange = (e) => {
    const value = e.target.type === 'number' ? parseFloat(e.target.value) : e.target.value;
    setFormData({ ...formData, [e.target.name]: value });
  };

  return (
    <form>
      {/* Your form inputs */}
      <input type="file" onChange={handleFileChange} />

        <input type="text" name="title" value={formData.title} onChange={handleChange} placeholder="Title" />
      <input type="text" name="type" value={formData.type} onChange={handleChange} placeholder="type" />
      <input type="text" name="publisher" value={formData.publisher} onChange={handleChange} placeholder="publisher" />

      <input type="number" name="score" value={formData.score} onChange={handleChange} placeholder="Score" />
      <input type="number" name="mal_id" value={formData.mal_id} onChange={handleChange} placeholder="malId" />
      <input type="number" name="popularity" value={formData.popularity} onChange={handleChange} placeholder="popularity" />

      <button type="button" onClick={handleUploadAndAddData}>Upload and Add Data</button>
    </form>
  );
}

export default AdminPage;