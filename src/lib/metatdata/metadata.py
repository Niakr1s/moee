import pathlib
import sys
from mutagen.id3 import ID3, USLT, TPE1, TIT2
from mutagen.oggvorbis import OggVorbis


# ----------- CLEAN ------------
def cleanMetadata(musicFilePath: str):
    print("start cleanMetadata to ", musicFilePath)
    ext = pathlib.Path(musicFilePath).suffix
    if ext == ".mp3":
        cleanMetadataMp3(musicFilePath)
    elif ext == ".ogg":
        cleanMetadataOgg(musicFilePath)
    else:
        raise Exception("wrong extension, want .mp3 or .ogg")
    print("end cleanMetadata to ", musicFilePath)


def cleanMetadataMp3(musicFilePath: str):
    print("start cleanMetadataMp3 to ", musicFilePath)
    audio = ID3(musicFilePath)
    audio.clear()
    audio.save()


def cleanMetadataOgg(musicFilePath: str):
    print("start cleanMetadataOgg to ", musicFilePath)
    audio = OggVorbis(musicFilePath)
    audio.clear()
    audio.save()


# ----------- LYRICS ------------
def writeLyrics(musicFilePath: str, contents: str):
    print("start writeLyrics to ", musicFilePath)
    ext = pathlib.Path(musicFilePath).suffix
    if ext == ".mp3":
        writeLyricsMp3(musicFilePath, contents)
    elif ext == ".ogg":
        writeLyricsOgg(musicFilePath, contents)
    else:
        raise Exception("wrong extension, want .mp3 or .ogg")
    print("end writeLyrics to ", musicFilePath)


def writeLyricsMp3(musicFilePath: str, contents: str):
    print("start writeLyricsMp3 to ", musicFilePath)
    audio = ID3(musicFilePath)
    audio.add(USLT(lang='   ', desc='', text=contents))
    audio.save()


def writeLyricsOgg(musicFilePath: str, contents: str):
    print("start writeLyricsOgg to ", musicFilePath)
    audio = OggVorbis(musicFilePath)
    audio['lyrics'] = contents
    audio.save()


# ----------- ARTIST ------------
def writeArtist(musicFilePath: str, contents: str):
    print("start writeArtist to ", musicFilePath)
    ext = pathlib.Path(musicFilePath).suffix
    if ext == ".mp3":
        writeArtistMp3(musicFilePath, contents)
    elif ext == ".ogg":
        writeArtistOgg(musicFilePath, contents)
    else:
        raise Exception("wrong extension, want .mp3 or .ogg")
    print("end writeArtist to ", musicFilePath)


def writeArtistMp3(musicFilePath: str, contents: str):
    print("start writeArtistMp3 to ", musicFilePath)
    audio = ID3(musicFilePath)
    audio.add(TPE1(text=contents))
    audio.save()


def writeArtistOgg(musicFilePath: str, contents: str):
    print("start writeArtistOgg to ", musicFilePath)
    audio = OggVorbis(musicFilePath)
    audio['artist'] = contents
    audio.save()


# ----------- TITLE ------------
def writeTitle(musicFilePath: str, contents: str):
    print("start writeTitle to ", musicFilePath)
    ext = pathlib.Path(musicFilePath).suffix
    if ext == ".mp3":
        writeTitleMp3(musicFilePath, contents)
    elif ext == ".ogg":
        writeTitleOgg(musicFilePath, contents)
    else:
        raise Exception("wrong extension, want .mp3 or .ogg")
    print("end writeTitle to ", musicFilePath)


def writeTitleMp3(musicFilePath: str, contents: str):
    print("start writeTitleMp3 to ", musicFilePath)
    audio = ID3(musicFilePath)
    audio.add(TIT2(text=contents))
    audio.save()


def writeTitleOgg(musicFilePath: str, contents: str):
    print("start writeTitleOgg to ", musicFilePath)
    audio = OggVorbis(musicFilePath)
    audio['title'] = contents
    audio.save()

# ----------------------------------


def run():
    usage = """Usage:
lyrics [filepath] [contents] - Writes lyrics to audiofile
artist [filepath] [contents] - Writes artist to audiofile
title [filepath] [contents] - Writes title to audiofile
clean [filepath] - Erases metadata from audiofile
"""

    # 0             1      2        3
    # "metadata.py" lyrics file.mp3 "some cool lyrics"
    args = sys.argv

    # 3 is minimum arguments
    if len(args) < 3:
        print(usage)
        raise Exception("not enough arguments")

    action = sys.argv[1]
    filePath = sys.argv[2]

    if action == "clean":
        cleanMetadata(filePath)
        return

    if len(args) < 4:
        print(usage)
        raise Exception("not enough arguments")
    contents = sys.argv[3]

    if action == "lyrics":
        writeLyrics(filePath, contents)
    elif action == "artist":
        writeArtist(filePath, contents)
    elif action == "title":
        writeTitle(filePath, contents)
    else:
        raise Exception("wrong action ", action)


if __name__ == "__main__":
    run()
