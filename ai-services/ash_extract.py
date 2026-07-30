import os
import subprocess
import tempfile

def extract_text(filepath: str) -> str:
    ext = os.path.splitext(filepath)[1].lower()
    
    if ext == '.pdf':
        return _extract_pdf(filepath)
    elif ext == '.docx':
        return _extract_docx(filepath)
    elif ext in ['.jpg', '.jpeg', '.png']:
        return _extract_image(filepath)
    elif ext in ['.mp3', '.wav', '.m4a']:
        return _extract_audio(filepath)
    elif ext in ['.zip', '.tar', '.gz']:
        return _extract_archive(filepath)
        
    return ""

def _extract_pdf(filepath: str) -> str:
    try:
        # try pdftotext (poppler)
        result = subprocess.run(['pdftotext', filepath, '-'], capture_output=True, text=True, check=True)
        return result.stdout
    except Exception:
        try:
            import PyPDF2
            text = ""
            with open(filepath, 'rb') as f:
                reader = PyPDF2.PdfReader(f)
                for page in reader.pages:
                    text += page.extract_text() + "\n"
            return text
        except Exception as e:
            return f"Error extracting PDF: {e}"

def _extract_docx(filepath: str) -> str:
    try:
        import docx
        doc = docx.Document(filepath)
        return "\n".join([para.text for para in doc.paragraphs])
    except Exception as e:
        return f"Error extracting DOCX: {e}"

def _extract_image(filepath: str) -> str:
    try:
        import pytesseract
        from PIL import Image
        return pytesseract.image_to_string(Image.open(filepath))
    except Exception as e:
        return f"Error extracting Image OCR: {e}"

def _extract_audio(filepath: str) -> str:
    try:
        import whisper
        model = whisper.load_model("base")
        result = model.transcribe(filepath)
        return result.get('text', '')
    except Exception as e:
        return f"Error extracting Audio: {e}"

def _extract_archive(filepath: str) -> str:
    try:
        if filepath.endswith('.zip'):
            import zipfile
            with zipfile.ZipFile(filepath, 'r') as z:
                return "\n".join(z.namelist())
        elif filepath.endswith('.tar.gz') or filepath.endswith('.tar'):
            import tarfile
            with tarfile.open(filepath, 'r') as t:
                return "\n".join(t.getnames())
    except Exception as e:
        return f"Error listing archive: {e}"
    return ""
