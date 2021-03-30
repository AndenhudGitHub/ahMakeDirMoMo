<?php
require 'JPEG_ICC.php';

$config    = file_get_contents("config.json");
$configArr = json_decode($config,true);

$smallPath    = file_get_contents("needToSize1000T1040.txt");
$smallPath    = str_replace("\n","",$smallPath);
$smallPathArr    = array_unique(explode(';',$smallPath));

$bigImgPath = $configArr['1000A1040Blank'];

if(count($smallPathArr)>0){
    foreach($smallPathArr as $filepath){
        if($filepath!="")
            mergePic($bigImgPath,$filepath);
    }
}

function mergePic($bigImgPath,$qCodePath){
    $MyJpeg   = new JPEG_ICC();
    $bigImg = imagecreatefromstring(file_get_contents($bigImgPath));
    $qCodeImg = imagecreatefromstring(file_get_contents($qCodePath));

    list($qCodeWidth, $qCodeHight, $qCodeType) = getimagesize($qCodePath);

    imagecopymerge($bigImg, $qCodeImg, 0, 20, 0, 0, $qCodeWidth, $qCodeHight, 100);
    
    list($bigWidth, $bigHight, $bigType) = getimagesize($bigImgPath);

    imagejpeg($bigImg,$qCodePath."-bigger",100);
    if ($MyJpeg->LoadFromJPEG($qCodePath)) {
        $MyJpeg->SaveToJPEG($qCodePath."-bigger");
    }
    rename($qCodePath."-bigger",$qCodePath);
    echo "改照片大小1000X1040，".$qCodePath."，成功\n";
}

