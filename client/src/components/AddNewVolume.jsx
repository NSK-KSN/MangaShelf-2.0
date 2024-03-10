import React, { useState, useEffect } from 'react';
import axios from 'axios';

function AddNewVolume() {
  const [formData, setFormData] = useState({
    title_id: '',
    volume_number: '',
    release_date: '',
    cover_link: '',
    imageFile: null // State to store the selected image file
  });

  const [mangaTitles, setMangaTitles] = useState([]);
  const [selectedFile, setSelectedFile] = useState(null);

  useEffect(() => {
    // Fetch manga titles from the backend
    const fetchMangaTitles = async () => {
      try {
        const response = await axios.get('http://localhost:8080/fetch-manga-titles');
        setMangaTitles(response.data);
      } catch (error) {
        console.error('Error fetching manga titles:', error);
      }
    };

    fetchMangaTitles();
  }, []);

  const handleFileChange = (e) => {
    setSelectedFile(e.target.files[0]);
  };

  const handleChange = (e) => {
    const value = e.target.type === 'number' ? parseFloat(e.target.value) : e.target.value;
    setFormData({ ...formData, [e.target.name]: value });
  };

  const handleUploadAndAddData = async () => {
    console.log(formData.title_id)
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
      cover_link: coverImageUrl, // Update cover_link data
      title: parseInt(formData.title),
    };
    try {
      const addDataResponse = await axios.post('http://localhost:8080/add-datas', formDataDB);
      console.log(addDataResponse.data);
    } catch (error) {
      console.error('Error adding data:', error);
    }
  };

  return (
    <div>
      <h1>Add new volume</h1>
      <form>
        <input type="file" onChange={handleFileChange} />

        <select name="title_id" value={formData.title_id} onChange={handleChange}>
          <option value="">Select Type</option>
            {mangaTitles.map(mangaTitle => (
              <option key={mangaTitle.id} value={mangaTitle.id}>{mangaTitle.name} ({mangaTitle.type})</option>
            ))}
        </select>

        <input type="number" name="volume_number" value={formData.volume_number} onChange={handleChange} placeholder="volume_number" />

        <input type="month" placeholder="Release Date" value={formData.release_date} onChange={handleChange} name="release_date" />

        <button type="button" onClick={handleUploadAndAddData}>Upload and Add Data</button>
      </form>
    </div>
  );
}

export default AddNewVolume;
