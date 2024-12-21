package com.codersee.springkotlindocker.controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestParam
import org.springframework.web.bind.annotation.RestController
import java.util.concurrent.atomic.AtomicInteger
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpStatus
import org.springframework.http.MediaType
import org.springframework.web.bind.annotation.*
import org.springframework.web.multipart.MultipartFile
import java.awt.Color
import java.awt.Graphics2D
import java.awt.image.BufferedImage
import javax.imageio.ImageIO
import javax.servlet.http.HttpServletResponse
import kotlinx.coroutines.*

@RestController
class Controller {

        @PostMapping("/convert")
        @ResponseStatus(HttpStatus.OK)
        fun convertToBlackAndWhite(
                @RequestParam("image") file: MultipartFile,
                response: HttpServletResponse
        ) = runBlocking {
                if (file != null && !file.isEmpty) {
                        val originalImage = ImageIO.read(file.inputStream)
                        val bwImage = convertImage(originalImage)

                        response.contentType = MediaType.IMAGE_PNG_VALUE
                        ImageIO.write(bwImage, "png", response.outputStream)
                        response.outputStream.flush()
                }
        }

        suspend fun convertImage(originalImage: BufferedImage): BufferedImage {
                val (topLeft, topRight, bottomLeft, bottomRight) = splitImage(originalImage)
                val processedImages = coroutineScope {
                val deferred = listOf(
                        async { convertImageToBlackAndWhite(topLeft) },
                        async { convertImageToBlackAndWhite(topRight) },
                        async { convertImageToBlackAndWhite(bottomLeft) },
                        async { convertImageToBlackAndWhite(bottomRight) }
                )
                deferred.awaitAll()
                }
                val outputImage = combineImages(processedImages)
                return outputImage
        }

        suspend fun convertImageToBlackAndWhite(originalImage: BufferedImage): BufferedImage {
                val width = originalImage.width
                val height = originalImage.height
                val outputImage = BufferedImage(width, height, BufferedImage.TYPE_BYTE_GRAY)
                val g: Graphics2D = outputImage.createGraphics()

                g.drawImage(originalImage, 0, 0, null)
                g.dispose()

                return outputImage
        }

        private fun splitImage(image: BufferedImage): List<BufferedImage> {
                val width = image.width
                val height = image.height
                val midWidth = width / 2
                val midHeight = height / 2

                return listOf(
                        image.getSubimage(0, 0, midWidth, midHeight),
                        image.getSubimage(midWidth, 0, midWidth, midHeight),
                        image.getSubimage(0, midHeight, midWidth, midHeight),
                        image.getSubimage(midWidth, midHeight, midWidth, midHeight)
                )
        }

        private fun combineImages(images: List<BufferedImage>): BufferedImage {
                val width = images[0].width * 2
                val height = images[0].height * 2
                val combinedImage = BufferedImage(width, height, BufferedImage.TYPE_INT_RGB)
        
                val graphics = combinedImage.createGraphics()
                graphics.drawImage(images[0], 0, 0, null)
                graphics.drawImage(images[1], images[0].width, 0, null)
                graphics.drawImage(images[2], 0, images[0].height, null)
                graphics.drawImage(images[3], images[0].width, images[0].height, null)
                graphics.dispose()
        
                return combinedImage
        }
}