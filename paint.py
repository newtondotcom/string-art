import numpy as np
from PIL import Image
import matplotlib.pyplot as plt
from matplotlib.backends.backend_pdf import PdfPages

REGION_SIZE = 5
IMAGE_PATH = 'images/meliodas3.jpg'
OUTPUT = 'JPG'  # 'PDF' or 'JPG'
MAX_DOT_SIZE = REGION_SIZE * 12 # this value changes depending on the region size and the format of the image

def load_image(image_path):
    """Load an image from a file path."""
    return Image.open(image_path)

def convert_to_grayscale(image):
    """Convert an image to grayscale."""
    return image.convert('L')

def image_to_varying_dots(image, region_size):
    """Convert a grayscale image to an array of dots with varying sizes."""
    image_array = np.array(image)
    dot_array = np.zeros((image_array.shape[0] // region_size, image_array.shape[1] // region_size))

    for i in range(dot_array.shape[0]):
        for j in range(dot_array.shape[1]):
            region = image_array[i*region_size:(i+1)*region_size, j*region_size:(j+1)*region_size]
            dot_array[i, j] = region.mean()

    return dot_array

def plot_and_save_varying_dots(dot_array, region_size, save_path, image_size, hollow=False, dpi=300, edge_width=0.1, output_format='PDF'):
    """Plot the array of dots with varying sizes and save as a PDF or JPG, centered on an A5 page."""
    a5_width, a5_height = 8.27, 5.83  # A5 dimensions in inches

    # Get the original image dimensions
    original_width, original_height = image_size

    # Scale the image to fit within A5 while maintaining aspect ratio
    aspect_ratio = original_width / original_height
    if aspect_ratio > (a5_width / a5_height):
        fig_width = a5_width
        fig_height = a5_width / aspect_ratio
    else:
        fig_height = a5_height
        fig_width = a5_height * aspect_ratio

    # Calculate offsets to center the image on the A5 page
    x_offset = (a5_width - fig_width) / 2
    y_offset = (a5_height - fig_height) / 2

    fig, ax = plt.subplots(figsize=(fig_width, fig_height), dpi=dpi)

    for i in range(dot_array.shape[0]):
        for j in range(dot_array.shape[1]):
            dot_size = (255 - dot_array[i, j]) / 255 * MAX_DOT_SIZE  # Invert brightness for dot size
            if hollow:
                ax.scatter(j * region_size + x_offset, i * region_size + y_offset, s=dot_size, edgecolors='black', facecolors='none', linewidth=edge_width)
            else:
                ax.scatter(j * region_size + x_offset, i * region_size + y_offset, s=dot_size, color='black', linewidth=edge_width)

    ax.invert_yaxis()
    ax.axis('off')
    plt.subplots_adjust(left=0, right=1, top=1, bottom=0)

    # Save the figure
    if output_format == 'PDF':
        # Save as PDF
        with PdfPages(save_path) as pdf:
            pdf.savefig(fig, bbox_inches='tight', pad_inches=0)
    elif output_format == 'JPG':
        # Save as JPG (to save in RGB mode)
        fig.savefig(save_path, format='jpg', dpi=dpi, bbox_inches='tight', pad_inches=0)

    plt.close(fig)

image = load_image(IMAGE_PATH)
image_name = image.filename
image_extension = image_name.split('.')[-1]
grayscale_image = convert_to_grayscale(image)

# Generate the dot array
dot_array = image_to_varying_dots(grayscale_image, REGION_SIZE)

# Get the original image size
original_size = grayscale_image.size

# Output paths based on the specified output format
if OUTPUT == 'PDF':
    filled_output_path = image_name.replace('.'+image_extension, '_filled_a5_centered.pdf')
    hollow_output_path = image_name.replace('.'+image_extension, '_hollow_a5_centered.pdf')
elif OUTPUT == 'JPG':
    filled_output_path = image_name.replace('.'+image_extension '_filled_a5_centered.jpg')
    hollow_output_path = image_name.replace('.'+image_extension, '_hollow_a5_centered.jpg')

# Plot the dots with varying sizes and save the filled image
plot_and_save_varying_dots(dot_array, REGION_SIZE, filled_output_path, original_size, edge_width=0.1, output_format=OUTPUT)

# Plot the dots with varying sizes and save the hollow image
plot_and_save_varying_dots(dot_array, REGION_SIZE, hollow_output_path, original_size, hollow=True, edge_width=0.1, output_format=OUTPUT)
